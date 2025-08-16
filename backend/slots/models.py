# slots/models.py
from datetime import timedelta, datetime
import uuid
from django.db import models
from django.conf import settings
from .encoders import PrettyJSONEncoder
from django.utils.translation import gettext_lazy as _


class SlotType(models.Model):
    """
    Different types of slots, e.g., 'event', 'document', 'facility'.
    Frequency is used to determine trending/popular slots.
    """
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    name = models.CharField(max_length=100, unique=True)
    frequency = models.PositiveIntegerField(default=0)  # trending score / usage count
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        ordering = ['-frequency', 'name']

    def __str__(self):
        return self.name


class Slot(models.Model):
    """
    Actual slot entry under a slot type.
    """
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    slot_type = models.ForeignKey(SlotType, on_delete=models.CASCADE, related_name="slots")
    title = models.CharField(max_length=255)
    description = models.JSONField(_("description"), blank=True, null=True, encoder=PrettyJSONEncoder)
    action = models.JSONField(_("action"), blank=True, null=True, encoder=PrettyJSONEncoder)
    fields = models.JSONField(_("fields"), blank=True, null=True, encoder=PrettyJSONEncoder)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    cover_image = models.ImageField(upload_to="slots/covers/", blank=True, null=True)

    owner = models.ForeignKey(
        settings.AUTH_USER_MODEL,
        on_delete=models.CASCADE,
        related_name="owned_slots"
    )

    collaborators = models.ManyToManyField(
        settings.AUTH_USER_MODEL,
        blank=True,
        related_name="collaborating_slots",
        help_text="Users who can edit/manage this slot."
    )

    class Meta:
        ordering = ['-created_at']

    def __str__(self):
        return self.title


class SlotSchedule(models.Model):
    slot = models.ForeignKey(Slot, on_delete=models.CASCADE, related_name="schedules")
    start_time = models.DateTimeField(blank=True, null=True)
    duration_minutes = models.PositiveIntegerField(default=20)
    grace_period_minutes = models.PositiveIntegerField(default=0)
    recurring = models.BooleanField(default=False)
    recurrence_pattern = models.JSONField(blank=True, null=True)
    flexible = models.BooleanField(default=False)
    price = models.DecimalField(max_digits=20, decimal_places=2, default=0)

    class Meta:
        ordering = ['start_time']

    def is_available(self, desired_time=None):
        """
        Check if schedule is available.
        - If flexible=True, desired_time must be provided
        - If fixed, just check the scheduled start_time
        """

        if self.flexible:
            if not desired_time:
                raise ValueError("Must provide desired_time for flexible booking")
            start = desired_time
        else:
            start = self.start_time

        end = start + timedelta(minutes=self.duration_minutes + self.grace_period_minutes)

        # Check for overlapping bookings
        overlapping = Booking.objects.filter(
            schedule__slot=self.slot,
            schedule__flexible=self.flexible,
            status__in=['pending', 'confirmed']
        ).filter(
            # Only for fixed slots, check overlap with schedule
            schedule__start_time__lt=end if not self.flexible else datetime.max,
            schedule__start_time__gte=start if not self.flexible else datetime.min
        )

        return not overlapping.exists()


class Booking(models.Model):
    STATUS_CHOICES = [
        ('pending', 'Pending'),
        ('confirmed', 'Confirmed'),
        ('cancelled', 'Cancelled'),
    ]

    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)

    wallet = models.ForeignKey("wallet.Wallet", on_delete=models.CASCADE, related_name="bookings")
    schedule = models.ForeignKey("SlotSchedule", on_delete=models.CASCADE, null=True, related_name="schedule_bookings")

    desired_time = models.DateTimeField(blank=True, null=True)  # only used for flexible slots

    booking_date = models.DateTimeField(auto_now_add=True)
    status = models.CharField(max_length=20, choices=STATUS_CHOICES, default='pending')
    details = models.JSONField(blank=True, null=True)
    payment_transaction = models.ForeignKey("wallet.Transaction", null=True, blank=True, on_delete=models.SET_NULL)

    class Meta:
        ordering = ['-booking_date']

    def __str__(self):
        return f"Booking {self.id} for {self.schedule.slot.title} at {self.schedule.start_time}"

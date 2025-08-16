from celery import shared_task
from django.utils import timezone
from django.core.mail import send_mail
from datetime import timedelta

from slots.models import Booking


@shared_task
def send_booking_reminders():
    """
    Send reminders for bookings happening 1 day from now.
    Only send if booking is confirmed and has a schedule.
    """
    tomorrow = timezone.now() + timedelta(days=1)
    bookings = Booking.objects.filter(
        status='confirmed',
        schedule__start_time__date=tomorrow.date()
    )

    for booking in bookings:
        user = booking.wallet.user
        schedule = booking.schedule

        if not user.email:
            continue  # skip if user has no email

        subject = "📅 Reminder: Your booking is tomorrow"
        message = (
            f"Hello {user.get_full_name() or user.username},\n\n"
            f"This is a reminder that you have a booking for:\n"
            f"🔹 {schedule.slot.title}\n"
            f"📍 Scheduled at: {schedule.start_time.strftime('%Y-%m-%d %H:%M')}\n\n"
            f"Thank you for using our service!"
        )

        send_mail(
            subject,
            message,
            None,  # uses DEFAULT_FROM_EMAIL from settings
            [user.email],
            fail_silently=False,
        )

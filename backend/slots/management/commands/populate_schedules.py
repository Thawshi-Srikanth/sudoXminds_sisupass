from datetime import datetime, timedelta
from django.core.management.base import BaseCommand
from slots.models import Slot, SlotSchedule


class Command(BaseCommand):
    help = "Populate SlotSchedule for existing Slots"

    def handle(self, *args, **options):
        slots = Slot.objects.all()
        now = datetime.now()

        for slot in slots:
            # Example: create 5 schedules per slot, 1-hour apart
            for i in range(5):
                start_time = now + timedelta(hours=i)
                duration = 20  # minutes
                grace = 5      # minutes
                recurring = False

                # Skip if schedule already exists for this slot at this time
                if SlotSchedule.objects.filter(slot=slot, start_time=start_time).exists():
                    continue

                SlotSchedule.objects.create(
                    slot=slot,
                    start_time=start_time,
                    duration_minutes=duration,
                    grace_period_minutes=grace,
                    recurring=recurring,
                    flexible=False
                )
                self.stdout.write(self.style.SUCCESS(
                    f"Created schedule for slot '{slot.title}' at {start_time}"
                ))

        self.stdout.write(self.style.SUCCESS("SlotSchedule population completed."))

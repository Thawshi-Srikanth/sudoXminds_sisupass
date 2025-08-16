from django.core.management.base import BaseCommand
from slots.models import SlotType


class Command(BaseCommand):
    help = "Populate the SlotType table with predefined categories"

    def handle(self, *args, **kwargs):
        categories = [
            {"name": "Events", "description": "Workshops, seminars, and student programs", "frequency": 10},
            {"name": "Facilities", "description": "Book labs, and sports areas", "frequency": 9},
            {"name": "Documents", "description": "Request IDs, transcripts, and certificates", "frequency": 8},
            {"name": "Appointments", "description": "Meet with counselors and staff", "frequency": 8},
            {"name": "Exams", "description": "Register for tests and practicals", "frequency": 5},
            {"name": "Health", "description": "Checkups, vaccinations, and counseling", "frequency": 4},
            {"name": "Transport", "description": "Book buses, campus shuttles, and permits", "frequency": 3},
            {"name": "Accommodation", "description": "Reserve hostel rooms and campus housing", "frequency": 2},
        ]

        for cat in categories:
            obj, created = SlotType.objects.get_or_create(
                name=cat["name"],
                defaults={"frequency": cat["frequency"]}
            )
            if created:
                self.stdout.write(self.style.SUCCESS(f"Added: {obj.name}"))
            else:
                self.stdout.write(self.style.WARNING(f"Skipped (exists): {obj.name}"))

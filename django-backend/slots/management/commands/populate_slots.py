import uuid
import random
from django.core.management.base import BaseCommand
from django.contrib.auth import get_user_model
from slots.models import SlotType, Slot

User = get_user_model()


class Command(BaseCommand):
    help = "Populate the database with creative sample slots (cover image only)"

    def handle(self, *args, **kwargs):
        users = list(User.objects.all())
        if not users:
            self.stdout.write(self.style.ERROR(
                "No users found! Create some users first."))
            return

        slot_types = SlotType.objects.all()
        if not slot_types:
            self.stdout.write(self.style.ERROR("No SlotTypes found!"))
            return

        self.stdout.write("Populating slots...")

        def generate_fields_for_slot(slot_type_name):
            """Generate slot-type-specific form fields with variation"""
            fields = []

            if slot_type_name.lower() == "documents":
                fields.append({
                    "title": "Personal Info",
                    "fields": [
                        {"name": "first_name", "label": "First Name",
                            "type": "text", "required": True},
                        {"name": "last_name", "label": "Last Name",
                            "type": "text", "required": True},
                        {"name": "date_of_birth", "label": "Date of Birth",
                            "type": "date", "required": True},
                        {"name": "nationality", "label": "Nationality",
                            "type": "text", "required": True},
                    ]
                })
                fields.append({
                    "title": "Document Details",
                    "fields": [
                        {"name": "document_type", "label": "Document Type", "type": "select", "options": [
                            "Passport", "NIC"], "required": True},
                        {"name": "upload_file", "label": "Upload File",
                            "type": "file", "required": True},
                    ]
                })
                if random.choice([True, False]):
                    fields.append({
                        "title": "Additional Notes",
                        "fields": [
                            {"name": "notes", "label": "Notes",
                                "type": "textarea", "required": False}
                        ]
                    })
            elif slot_type_name.lower() == "events":
                fields.append({
                    "title": "Event Registration",
                    "fields": [
                        {"name": "full_name", "label": "Full Name",
                            "type": "text", "required": True},
                        {"name": "attendees", "label": "Number of Attendees",
                            "type": "number", "required": False},
                        {"name": "session", "label": "Preferred Session", "type": "select", "options": [
                            "Morning", "Afternoon", "Evening"], "required": False}
                    ]
                })
                if random.choice([True, False]):
                    fields.append({
                        "title": "Optional Notes",
                        "fields": [
                            {"name": "notes", "label": "Notes",
                                "type": "textarea", "required": False}
                        ]
                    })
            elif slot_type_name.lower() == "facilities":
                fields.append({
                    "title": "Booking Details",
                    "fields": [
                        {"name": "date", "label": "Booking Date",
                            "type": "date", "required": True},
                        {"name": "duration_hours",
                            "label": "Duration (hours)", "type": "number", "required": True}
                    ]
                })
                if random.choice([True, False]):
                    fields.append({
                        "title": "Equipment",
                        "fields": [
                            {"name": "projector_needed", "label": "Projector Needed",
                                "type": "checkbox", "required": False}
                        ]
                    })
            else:
                fields.append({
                    "title": "Basic Info",
                    "fields": [
                        {"name": "name", "label": "Name",
                            "type": "text", "required": True}
                    ]
                })
            return fields

        # Now create slots
        for slot_type in slot_types:
            slots_to_create = []

            for i in range(1, 4):  # 3 slots per type
                owner = random.choice(users)
                collaborators = random.sample(
                    users, k=random.randint(0, min(3, len(users))))

                # Dynamic description
                description = {
                    "about": f"{slot_type.name} Slot {i} - creatively generated description",
                    "location": {
                        "address": f"{random.randint(1, 100)} Main Street, Colombo",
                        "latitude": 6.9271 + random.uniform(-0.01, 0.01),
                        "longitude": 79.8612 + random.uniform(-0.01, 0.01)
                    },
                    "event": {
                        "start_date": f"2025-09-{10+i}",
                        "end_date": f"2025-09-{10+i}",
                        "time": f"{9+i}:00 - {11+i}:00"
                    },
                    "schedules": {
                        "passport_booking": "Available every Monday",
                        "nic_booking": "Available Wednesdays"
                    },
                    "facilities": random.sample(["Free parking", "Wheelchair access", "On-site printing", "Cafeteria", "Projector"], k=3)
                }

                # Random action type
                action_type = random.choice(['form', 'link', 'none'])
                if action_type == 'form':
                    action = {
                        "type": "form",
                        "label": "Apply Now",
                        "form_fields": generate_fields_for_slot(slot_type.name)
                    }
                elif action_type == 'link':
                    action = {
                        "type": "link",
                        "label": "Book Online",
                        "url": "https://example.com/booking"
                    }
                else:
                    action = None

                slot = Slot(
                    slot_type=slot_type,
                    title=f"{slot_type.name} Example Slot {i}",
                    description=description,
                    action=action,
                    fields=action.get("form_fields") if action and action.get(
                        "type") == "form" else None,
                    owner=owner,
                    cover_image=f"https://picsum.photos/seed/{uuid.uuid4()}/800/400"
                )
                slots_to_create.append(slot)

            # Bulk create
            Slot.objects.bulk_create(slots_to_create)

            # Assign collaborators after creation
            for slot in Slot.objects.filter(slot_type=slot_type).order_by('-created_at')[:3]:
                slot.collaborators.set(random.sample(
                    users, k=random.randint(0, min(3, len(users)))))
                slot.save()

        self.stdout.write(self.style.SUCCESS(
            "Slots populated successfully with cover images, varied forms, and actions!"))

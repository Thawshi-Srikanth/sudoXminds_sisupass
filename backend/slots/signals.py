from django_celery_beat.models import PeriodicTask, IntervalSchedule
import json


def create_booking_reminder_task(sender, **kwargs):
    try:
        schedule, created = IntervalSchedule.objects.get_or_create(
            every=1,
            period=IntervalSchedule.DAYS
        )

        PeriodicTask.objects.get_or_create(
            interval=schedule,
            name='Send booking reminders',
            task='slots.tasks.send_booking_reminders',
            defaults={
                'args': json.dumps([]),
            }
        )
    except Exception as e:
        # Ignore errors if DB isn't ready
        print("Could not create periodic task:", e)

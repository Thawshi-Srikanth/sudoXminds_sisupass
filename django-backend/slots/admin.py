# slots/admin.py
from django.contrib import admin
from unfold.admin import ModelAdmin
from .models import SlotType, Slot, Booking, SlotSchedule


@admin.register(SlotType)
class SlotTypeAdmin(ModelAdmin):
    list_display = ['name', 'frequency', 'created_at']


@admin.register(Slot)
class SlotAdmin(ModelAdmin):
    list_display = ('title', 'slot_type', 'owner', 'created_at')
    search_fields = ('title',)
    list_filter = ('slot_type', 'created_at')

    # This makes the ManyToMany field appear as a dual list widget
    filter_horizontal = ('collaborators',)

    #readonly_fields = ["description", "action", "fields"]


@admin.register(Booking)
class BookingAdmin(ModelAdmin):
    list_display = ['id', 'wallet', 'status', 'booking_date']
    list_filter = ['status']


@admin.register(SlotSchedule)
class SlotScheduleAdmin(ModelAdmin):
    list_display = ['id',  'start_time', 'duration_minutes', 'grace_period_minutes']
    list_filter = ['slot', 'start_time', 'recurring', 'flexible']

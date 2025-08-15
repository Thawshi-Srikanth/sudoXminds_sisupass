# slots/admin.py
from django.contrib import admin
from .models import SlotType, Slot, Booking

@admin.register(SlotType)
class SlotTypeAdmin(admin.ModelAdmin):
    list_display = ['name', 'frequency', 'created_at']


@admin.register(Slot)
class SlotAdmin(admin.ModelAdmin):
    list_display = ['title', 'slot_type', 'created_at']
    list_filter = ['slot_type']


@admin.register(Booking)
class BookingAdmin(admin.ModelAdmin):
    list_display = ['id', 'slot', 'wallet', 'status', 'booking_date']
    list_filter = ['status']

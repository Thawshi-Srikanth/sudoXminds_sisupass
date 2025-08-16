from django.contrib import admin
from unfold.admin import ModelAdmin

from django.contrib.auth import get_user_model

CustomUser = get_user_model()


@admin.register(CustomUser)
class CustomUserAdmin(ModelAdmin):
    list_display = ['email', 'first_name', 'last_name', 'is_active', 'is_staff']
    search_fields = ['email', 'first_name', 'last_name']

from django.contrib import admin
from unfold.admin import ModelAdmin

from wallet.models import Wallet, Transaction, LocationType, Location, PassCategory, PassDetails


@admin.register(Wallet)
class WalletAdmin(ModelAdmin):
    list_display = ['user', 'balance', 'created_at']
    search_fields = ['user__email']


@admin.register(Transaction)
class TransactionAdmin(ModelAdmin):
    list_display = ['transaction_id', 'amount', 'transaction_date']
    list_filter = ['transaction_date']


@admin.register(LocationType)
class LocationTypeAdmin(ModelAdmin):
    list_display = ['name']


@admin.register(Location)
class LocationAdmin(ModelAdmin):
    list_display = ['name', 'location_type']


@admin.register(PassCategory)
class PassCategoryAdmin(ModelAdmin):
    list_display = ['name']


@admin.register(PassDetails)
class PassDetailsAdmin(ModelAdmin):
    list_display = ['wallet', 'category']

# wallet/schema/types.py
import graphene
from graphene_django import DjangoObjectType
from ..models import Wallet, Transaction, LocationType, Location, PassCategory, PassDetails


class WalletType(DjangoObjectType):
    class Meta:
        model = Wallet
        fields = ("wallet_id", "balance", "wallet_type",
                  "created_at", "updated_at")


class TransactionType(DjangoObjectType):
    class Meta:
        model = Transaction
        fields = ("transaction_id", "from_wallet", "to_wallet", "amount",
                  "transaction_type", "status", "transaction_date", "description")


class LocationTypeType(DjangoObjectType):
    class Meta:
        model = LocationType
        fields = ("id", "name", "description")


class LocationsType(DjangoObjectType):
    class Meta:
        model = Location
        fields = ("id", "name", "address", "latitude",
                  "longitude", "location_type")


class PassCategoryType(DjangoObjectType):
    class Meta:
        model = PassCategory
        fields = ("id", "name", "description", "allowed_location_types")


class PassDetailsType(DjangoObjectType):
    class Meta:
        model = PassDetails
        fields = (
            "id",
            "wallet",
            "category",
            "start_date",
            "end_date",
            "from_location",
            "to_location",
            "allowed_locations",
        )

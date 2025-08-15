import graphene
from graphene_django import DjangoObjectType
from slots.models import SlotType, Slot, Booking


class SlotTypeType(DjangoObjectType):
    class Meta:
        model = SlotType
        fields = "__all__"


class SlotTypeSlot(DjangoObjectType):
    class Meta:
        model = Slot
        fields = "__all__"


class SlotTypeBooking(DjangoObjectType):
    class Meta:
        model = Booking
        fields = "__all__"

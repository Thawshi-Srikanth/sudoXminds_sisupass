from graphene_django import DjangoObjectType
from slots.models import SlotType, Slot, Booking


class SlotTypeNode(DjangoObjectType):
    class Meta:
        model = SlotType
        fields = "__all__"


class SlotNode(DjangoObjectType):
    class Meta:
        model = Slot
        fields = "__all__"


class BookingNode(DjangoObjectType):
    class Meta:
        model = Booking
        fields = "__all__"

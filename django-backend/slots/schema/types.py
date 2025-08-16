import graphene
from graphene_django import DjangoObjectType
from slots.models import SlotSchedule, SlotType, Slot, Booking


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


class SlotScheduleType(DjangoObjectType):
    available = graphene.Boolean()

    class Meta:
        model = SlotSchedule
        fields = "__all__"

    def resolve_available(self, info):
        return self.is_available()

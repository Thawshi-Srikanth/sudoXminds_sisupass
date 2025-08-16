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
    slot_name = graphene.String()
    start_time = graphene.String()

    class Meta:
        model = Booking
        fields = "__all__"

    def resolve_slot_name(self, info):
        # Access the related slot title
        return self.schedule.slot.title if self.schedule else None

    def resolve_slot_time(self, info):
        # Access the related slot time
        return self.schedule.start_time if self.schedule else None


class SlotScheduleType(DjangoObjectType):
    available = graphene.Boolean()

    class Meta:
        model = SlotSchedule
        fields = "__all__"

    def resolve_available(self, info):
        return self.is_available()

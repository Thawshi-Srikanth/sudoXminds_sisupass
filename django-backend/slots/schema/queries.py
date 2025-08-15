import graphene
from .types import SlotTypeType, SlotTypeSlot, SlotTypeBooking
from slots.models import SlotType, Slot, Booking


class SlotTypeQuery(graphene.ObjectType):
    slot_types = graphene.List(SlotTypeType, trending=graphene.Boolean())
    slots = graphene.List(
        SlotTypeSlot, slot_type_id=graphene.UUID(required=False))
    bookings = graphene.List(
        SlotTypeBooking, wallet_id=graphene.UUID(required=False))

    def resolve_slot_types(self, info, trending=False):
        qs = SlotType.objects.all()
        if trending:
            qs = qs.order_by("-frequency")
        return qs

    def resolve_slots(self, info, slot_type_id=None):
        qs = Slot.objects.all()
        if slot_type_id:
            qs = qs.filter(slot_type_id=slot_type_id)
        return qs

    def resolve_bookings(self, info, wallet_id=None):
        qs = Booking.objects.all()
        if wallet_id:
            qs = qs.filter(wallet_id=wallet_id)
        return qs

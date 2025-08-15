import graphene
from .types import SlotTypeType, SlotTypeSlot, SlotTypeBooking
from slots.models import SlotType, Slot, Booking


class CreateSlotType(graphene.Mutation):
    class Arguments:
        name = graphene.String(required=True)
        frequency = graphene.Int(required=False)

    slot_type = graphene.Field(SlotTypeType)

    def mutate(self, info, name, frequency=0):
        obj = SlotType.objects.create(name=name, frequency=frequency)
        return CreateSlotType(slot_type=obj)


class CreateSlot(graphene.Mutation):
    class Arguments:
        slot_type_id = graphene.UUID(required=True)
        title = graphene.String(required=True)
        description = graphene.JSONString(required=False)
        action = graphene.JSONString(required=False)
        fields = graphene.JSONString(required=False)

    slot = graphene.Field(SlotTypeSlot)

    def mutate(self, info, slot_type_id, title, description=None, action=None, fields=None):
        slot = Slot.objects.create(
            slot_type_id=slot_type_id,
            title=title,
            description=description,
            action=action,
            fields=fields
        )
        return CreateSlot(slot=slot)


class CreateBooking(graphene.Mutation):
    class Arguments:
        wallet_id = graphene.UUID(required=True)
        slot_id = graphene.UUID(required=True)
        details = graphene.JSONString(required=False)

    booking = graphene.Field(SlotTypeBooking)

    def mutate(self, info, wallet_id, slot_id, details=None):
        booking = Booking.objects.create(
            wallet_id=wallet_id, slot_id=slot_id, details=details)
        return CreateBooking(booking=booking)

import graphene
from .queries import SlotTypeQuery, SlotQuery
from .mutations import CreateSlotType, CreateSlot, CreateBooking


class Query(SlotTypeQuery, SlotQuery, graphene.ObjectType):
    pass


class Mutation(graphene.ObjectType):
    create_slot_type = CreateSlotType.Field()
    create_slot = CreateSlot.Field()
    create_booking = CreateBooking.Field()

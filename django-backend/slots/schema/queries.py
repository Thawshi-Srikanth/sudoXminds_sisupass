import graphene
from .types import SlotNode, BookingNode, SlotScheduleType, SlotTypeNode
from slots.models import SlotType, Slot, Booking


class SlotTypeQuery(graphene.ObjectType):
    slot_types = graphene.List(SlotTypeNode, trending=graphene.Boolean())
    slots = graphene.List(
        SlotNode, slot_type_id=graphene.UUID(required=False))
    bookings = graphene.List(
        BookingNode, wallet_id=graphene.UUID(required=False))

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


class SlotQuery(graphene.ObjectType):
    slots_by_type = graphene.List(
        SlotNode,
        type_id=graphene.UUID(required=True),
        search=graphene.String(required=False)
    )
    trending_slots = graphene.List(
        SlotNode,
        type_id=graphene.UUID(required=True),
        limit=graphene.Int(required=False)
    )
    upcoming_bookings = graphene.List(
        BookingNode
    )

    slot_schedules = graphene.List(SlotScheduleType, slot_id=graphene.ID(required=True))

    slot_by_id = graphene.Field(SlotNode, id=graphene.UUID(required=True))

    def resolve_slots_by_type(root, info, type_id, search=None):
        qs = Slot.objects.filter(slot_type_id=type_id)
        if search:
            qs = qs.filter(title__icontains=search)
        return qs.order_by("description__event__start_date")

    def resolve_trending_slots(root, info, type_id, limit=5):
        return Slot.objects.filter(slot_type_id=type_id).order_by("-updated_at")[:limit]

    def resolve_upcoming_bookings(root, info):
        user = info.context.user
        if user.is_anonymous:
            return Booking.objects.none()
        return Booking.objects.filter(wallet__user=user).order_by("booking_date")

    def resolve_slot_by_id(root, info, id):
        return Slot.objects.filter(id=id).first()

    def resolve_slot_schedules(self, info, slot_id):
        try:
            slot = Slot.objects.get(id=slot_id)
            return slot.schedules.all()
        except Slot.DoesNotExist:
            return []

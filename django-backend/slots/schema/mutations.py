import graphene
from .types import BookingNode, SlotTypeNode
from slots.models import SlotSchedule, SlotType, Slot, Booking
from wallet.models import Wallet, Transaction


class CreateSlotType(graphene.Mutation):
    class Arguments:
        name = graphene.String(required=True)
        frequency = graphene.Int(required=False)

    slot_type = graphene.Field(SlotTypeNode)

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

    slot = graphene.Field(SlotTypeNode)

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
    booking = graphene.Field(BookingNode)

    class Arguments:
        schedule_id = graphene.ID(required=True)
        desired_time = graphene.DateTime(required=False)
        # no need for wallet_id; payment goes to owner
        # optional: user can still select their wallet if needed
        user_wallet_id = graphene.ID(required=False)
        details = graphene.JSONString(required=False)

    def mutate(self, info, schedule_id, desired_time=None, user_wallet_id=None, details=None):
        user = info.context.user
        if not user.is_authenticated:
            raise Exception("Authentication required")

        # Get schedule and slot owner
        try:
            schedule = SlotSchedule.objects.get(id=schedule_id)
        except SlotSchedule.DoesNotExist:
            raise Exception("Schedule not found")

        slot_owner = schedule.slot.owner
        try:
            owner_wallet = Wallet.objects.get(user=slot_owner, wallet_type='main')
        except Wallet.DoesNotExist:
            raise Exception("Slot owner does not have a main wallet")

        # User wallet (optional, default to user's main wallet)
        if user_wallet_id:
            user_wallet = Wallet.objects.get(wallet_id=user_wallet_id, user=user)
        else:
            user_wallet = Wallet.objects.get(user=user, wallet_type='main')

        # Check availability
        if not schedule.is_available(desired_time=desired_time):
            raise Exception("Schedule not available")

        # Handle payment if price > 0
        transaction = None
        if schedule.price > 0:
            if user_wallet.balance < schedule.price:
                raise Exception("Insufficient balance in your wallet")

            # Deduct from user wallet
            user_wallet.spend(schedule.price)

            # Credit owner wallet
            owner_wallet.balance += schedule.price
            owner_wallet.save()

            # Create transaction record
            transaction = Transaction.objects.create(
                from_wallet=user_wallet,
                to_wallet=owner_wallet,
                amount=schedule.price,
                transaction_type='spend',
                status='completed',
                description=f"Booking payment for {schedule.slot.title}"
            )

        # Create booking
        booking = Booking.objects.create(
            wallet=user_wallet,  # booking is associated with user's wallet
            schedule=schedule,
            desired_time=desired_time if schedule.flexible else None,
            status='confirmed' if schedule.price == 0 else 'pending',
            payment_transaction=transaction,
            details={"price": float(schedule.price), **(details or {})}
        )

        return CreateBooking(booking=booking)

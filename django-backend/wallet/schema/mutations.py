# wallet/schema/mutations.py
from decimal import Decimal
from graphql import GraphQLError
from django.utils import timezone
import graphene
from graphql_jwt.decorators import login_required
from .types import LocationsType, LocationTypeType, PassCategoryType, PassDetailsType, WalletType, TransactionType
from ..models import Location, PassCategory, PassDetails, Transaction, Wallet, LocationType
from django.db import transaction as db_transaction


class AddFunds(graphene.Mutation):
    class Arguments:
        wallet_id = graphene.UUID(required=False)
        amount = graphene.Decimal(required=True)

    wallet = graphene.Field(WalletType)

    @login_required
    def mutate(self, info, wallet_id, amount):
        user = info.context.user
        if wallet_id is None:
            wallet = Wallet.objects.filter(wallet_id=wallet_id, user=user).first()
        else:
            wallet = Wallet.objects.filter(user=user, wallet_type='main').first()
        if not wallet:
            raise Exception("Wallet not found")
        if amount <= 0:
            raise Exception("Amount must be positive")
        wallet.balance += amount
        wallet.save()
        return AddFunds(wallet=wallet)


class SpendFunds(graphene.Mutation):
    class Arguments:
        wallet_id = graphene.UUID(required=True)
        amount = graphene.Decimal(required=True)

    wallet = graphene.Field(WalletType)

    @login_required
    def mutate(self, info, wallet_id, amount):
        user = info.context.user
        wallet = Wallet.objects.filter(wallet_id=wallet_id, user=user).first()
        if not wallet:
            raise Exception("Wallet not found")
        if amount <= 0:
            raise Exception("Amount must be positive")
        if amount > wallet.balance:
            raise Exception("Insufficient balance")
        wallet.balance -= amount
        wallet.save()
        return SpendFunds(wallet=wallet)


class CreatePassWallet(graphene.Mutation):
    wallet = graphene.Field(WalletType)

    @login_required
    def mutate(self, info):
        user = info.context.user

        # Find the main wallet for this user
        parent_wallet = Wallet.objects.filter(
            user=user, wallet_type='main').first()
        if not parent_wallet:
            raise Exception(
                "Main wallet not found. You must have a main wallet before creating a pass.")

        # Create pass wallet
        wallet = Wallet.objects.create(
            user=user,
            wallet_type='pass',
            parent_wallet=parent_wallet
        )
        return CreatePassWallet(wallet=wallet)


class WalletMutation(graphene.ObjectType):
    add_funds = AddFunds.Field()
    spend_funds = SpendFunds.Field()
    create_pass_wallet = CreatePassWallet.Field()


class CreateLocationType(graphene.Mutation):
    class Arguments:
        name = graphene.String(required=True)
        description = graphene.String()

    location_type = graphene.Field(LocationTypeType)

    @login_required
    def mutate(self, info, name, description=None):
        obj = LocationType.objects.create(
            name=name, description=description or "")
        return CreateLocationType(location_type=obj)


class CreateLocation(graphene.Mutation):
    class Arguments:
        name = graphene.String(required=True)
        address = graphene.String(required=True)
        latitude = graphene.Float(required=True)
        longitude = graphene.Float(required=True)
        location_type_id = graphene.Int(required=True)

    location = graphene.Field(LocationsType)

    @login_required
    def mutate(self, info, name, address, latitude, longitude, location_type_id):
        location_type = LocationsType.objects.get(pk=location_type_id)
        obj = Location.objects.create(
            name=name,
            address=address,
            latitude=latitude,
            longitude=longitude,
            location_type=location_type
        )
        return CreateLocation(location=obj)


class CreatePassCategory(graphene.Mutation):
    class Arguments:
        name = graphene.String(required=True)
        description = graphene.String()
        allowed_location_type_ids = graphene.List(graphene.Int)

    category = graphene.Field(PassCategoryType)

    @login_required
    def mutate(self, info, name, description=None, allowed_location_type_ids=None):
        obj = PassCategory.objects.create(
            name=name, description=description or "")
        if allowed_location_type_ids:
            obj.allowed_location_types.set(
                LocationType.objects.filter(id__in=allowed_location_type_ids))
        return CreatePassCategory(category=obj)


class CreatePass(graphene.Mutation):
    class Arguments:
        category_id = graphene.Int(required=True)
        from_location_id = graphene.Int()
        to_location_id = graphene.Int()
        allowed_location_ids = graphene.List(graphene.Int)
        start_date = graphene.Date()
        end_date = graphene.Date()

    pass_details = graphene.Field(PassDetailsType)

    @login_required
    def mutate(
        self, info, category_id, from_location_id=None, to_location_id=None,
        allowed_location_ids=None, start_date=None, end_date=None
    ):
        user = info.context.user
        wallet = Wallet.objects.filter(user=user).first()
        if not wallet:
            raise Exception("User has no wallet")

        category = PassCategory.objects.get(pk=category_id)
        pass_obj = PassDetails.objects.create(
            wallet=wallet,
            category=category,
            from_location=Location.objects.get(
                pk=from_location_id) if from_location_id else None,
            to_location=Location.objects.get(
                pk=to_location_id) if to_location_id else None
        )
        if allowed_location_ids:
            pass_obj.allowed_locations.set(
                Location.objects.filter(id__in=allowed_location_ids))

        return CreatePass(pass_details=pass_obj)


class PassMutation(graphene.ObjectType):
    create_location_type = CreateLocationType.Field()
    create_location = CreateLocation.Field()
    create_pass_category = CreatePassCategory.Field()
    create_pass = CreatePass.Field()


class PassDetailInput(graphene.InputObjectType):
    category_id = graphene.ID(required=True)
    from_location_id = graphene.ID()
    to_location_id = graphene.ID()
    allowed_location_ids = graphene.List(graphene.ID)
    start_date = graphene.types.datetime.Date()
    end_date = graphene.types.datetime.Date()


class CreatePassWalletWithDetails(graphene.Mutation):
    class Arguments:
        pass_details = graphene.List(PassDetailInput, required=True)
        name = graphene.String(required=True)
        exp_date = graphene.String(required=False)

    wallet = graphene.Field(lambda: WalletType)

    created_pass_details = graphene.List(lambda: PassDetailsType)

    @classmethod
    @db_transaction.atomic
    def mutate(cls, root, info, pass_details, name, exp_date=None):
        user = info.context.user
        if user.is_anonymous:
            raise GraphQLError("Authentication required")

        # Get user's main wallet
        try:
            parent_wallet = Wallet.objects.get(user=user, wallet_type="main")
        except Wallet.DoesNotExist:
            raise GraphQLError("User does not have a main wallet")

        # Always create a new "pass" wallet
        wallet = Wallet.objects.create(
            name=name,
            exp_date=exp_date,
            user=user,
            wallet_type="pass",
            parent_wallet=parent_wallet,
            balance=Decimal("0.0")
        )

        created_passes = []

        for pd in pass_details:
            try:
                category = PassCategory.objects.get(id=pd.category_id)
            except PassCategory.DoesNotExist:
                raise GraphQLError(f"Invalid category ID: {pd.category_id}")

            pass_detail = PassDetails.objects.create(
                wallet=wallet,
                category=category,
                start_date=pd.start_date or timezone.now().date(),
                end_date=pd.end_date,
                from_location=Location.objects.filter(
                    id=pd.from_location_id).first() if pd.from_location_id else None,
                to_location=Location.objects.filter(
                    id=pd.to_location_id).first() if pd.to_location_id else None,
            )

            if pd.allowed_location_ids:
                allowed_locs = Location.objects.filter(
                    id__in=pd.allowed_location_ids)
                pass_detail.allowed_locations.set(allowed_locs)

            pass_detail.save()
            created_passes.append(pass_detail)

        # Return the newly created wallet and passes
        return CreatePassWalletWithDetails(wallet=wallet, created_pass_details=created_passes)


class PassDetailMutation(graphene.ObjectType):
    create_pass_wallet_with_details = CreatePassWalletWithDetails.Field()


class CreateTransaction(graphene.Mutation):
    transaction = graphene.Field(TransactionType)

    class Arguments:
        from_wallet_id = graphene.ID(required=False)
        to_wallet_id = graphene.ID(required=False)
        amount = graphene.Decimal(required=True)
        transaction_type = graphene.String(required=True)  # topup, spend, transfer
        description = graphene.String(required=False)

    @classmethod
    @db_transaction.atomic
    def mutate(cls, root, info, amount, transaction_type, from_wallet_id=None, to_wallet_id=None, description=None):
        user = info.context.user
        if user.is_anonymous:
            raise GraphQLError("Authentication required")

        # Resolve from_wallet
        from_wallet = None
        if from_wallet_id:
            from_wallet = Wallet.objects.filter(pk=from_wallet_id, user=user).first()
            if not from_wallet:
                raise GraphQLError("From wallet not found for this user")
        elif transaction_type in ["spend", "transfer"]:
            from_wallet = Wallet.objects.filter(user=user, wallet_type="main").first()
            if not from_wallet:
                raise GraphQLError("Main wallet not found for this user")

        # Resolve to_wallet
        to_wallet = None
        if to_wallet_id:
            to_wallet = Wallet.objects.filter(pk=to_wallet_id).first()
            if not to_wallet:
                raise GraphQLError("To wallet not found")
        elif transaction_type == "topup":
            to_wallet = Wallet.objects.filter(user=user, wallet_type="main").first()
            if not to_wallet:
                raise GraphQLError("Main wallet not found for this user")

        # Create base transaction
        transaction = Transaction.objects.create(
            from_wallet=from_wallet,
            to_wallet=to_wallet,
            amount=amount,
            transaction_type=transaction_type,
            description=description,
            status="pending"
        )

        try:
            if transaction_type == "topup":
                if not to_wallet:
                    raise GraphQLError("Target wallet required for topup")
                to_wallet.topup(amount)

            elif transaction_type == "spend":
                if not from_wallet:
                    raise GraphQLError("Source wallet required for spend")
                from_wallet.spend(amount)

            elif transaction_type == "transfer":
                if not from_wallet or not to_wallet:
                    raise GraphQLError("Both from_wallet and to_wallet required for transfer")
                # Deduct first
                from_wallet.spend(amount)
                # Then add to target
                to_wallet.topup(amount)

            else:
                raise GraphQLError("Invalid transaction type")

            transaction.status = "completed"

        except Exception as e:
            transaction.status = "failed"
            transaction.save()
            raise GraphQLError(str(e))

        transaction.save()
        return CreateTransaction(transaction=transaction)


class TransactionMutation(graphene.ObjectType):
    create_transaction = CreateTransaction.Field()

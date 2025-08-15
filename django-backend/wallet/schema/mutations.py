# wallet/schema/mutations.py
from datetime import date
import graphene
from graphql_jwt.decorators import login_required
from .types import LocationsType, LocationTypeType, PassCategoryType, PassDetailsType, WalletType
from ..models import Location, PassCategory, PassDetails, Wallet, LocationType


class AddFunds(graphene.Mutation):
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

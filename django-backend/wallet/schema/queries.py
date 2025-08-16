# wallet/schema/queries.py
import graphene
from graphql_jwt.decorators import login_required
from .types import LocationsType, LocationTypeType, PassCategoryType, PassDetailsType, WalletType, TransactionType
from ..models import Location, PassCategory, PassDetails, Wallet, Transaction, LocationType


class WalletQuery(graphene.ObjectType):
    my_wallets = graphene.List(WalletType)
    my_transactions = graphene.List(TransactionType)
    my_pass_wallets = graphene.List(WalletType)

    @login_required
    def resolve_my_wallets(root, info):
        user = info.context.user
        return Wallet.objects.filter(user=user)

    @login_required
    def resolve_my_transactions(root, info):
        user = info.context.user
        return Transaction.objects.filter(from_wallet__user=user) | Transaction.objects.filter(to_wallet__user=user)

    @login_required
    def resolve_my_pass_wallets(root, info):
        user = info.context.user
        main_wallet = Wallet.objects.filter(
            user=user, wallet_type='main').first()
        if not main_wallet:
            return []
        return Wallet.objects.filter(parent_wallet=main_wallet, wallet_type='pass')


class PassQuery(graphene.ObjectType):
    all_location_types = graphene.List(LocationTypeType)
    all_locations = graphene.List(LocationsType)
    all_pass_categories = graphene.List(PassCategoryType)
    my_passes = graphene.List(PassDetailsType)

    @login_required
    def resolve_all_location_types(root, info):
        return LocationType.objects.all()

    @login_required
    def resolve_all_locations(root, info):
        return Location.objects.select_related("location_type").all()

    @login_required
    def resolve_all_pass_categories(root, info):
        return PassCategory.objects.prefetch_related("allowed_location_types").all()

    @login_required
    def resolve_my_passes(root, info):
        user = info.context.user
        return PassDetails.objects.select_related(
            "wallet", "category", "from_location", "to_location"
        ).prefetch_related(
            "allowed_locations"
        ).filter(wallet__user=user)


class TransactionQuery(graphene.ObjectType):
    user_transactions = graphene.List(
        TransactionType,
        limit=graphene.Int(required=False)
    )

    @login_required
    def resolve_user_transactions(self, info, limit=None):
        user = info.context.user

        # Step 1: Get main wallets of user
        user_wallets = user.wallets.all()

        # Step 2: Get child wallets (where parent_wallet belongs to user)
        child_wallets = Wallet.objects.filter(parent_wallet__in=user_wallets)

        # Step 3: Combine all wallet IDs
        all_wallet_ids = list(user_wallets.values_list('wallet_id', flat=True)) + \
            list(child_wallets.values_list('wallet_id', flat=True))

        # Step 4: Get transactions where from_wallet OR to_wallet belongs to any of these wallets
        transactions = Transaction.objects.filter(
            from_wallet_id__in=all_wallet_ids
        ).order_by('-transaction_date')

        if limit:
            transactions = transactions[:limit]

        return transactions

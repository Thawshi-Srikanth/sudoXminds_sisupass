import graphene
from graphql import GraphQLError
from wallet.models import Wallet, Transaction
from .queries import WalletType
from wallet.utils import login_required


class TopUpWallet(graphene.Mutation):
    class Arguments:
        wallet_id = graphene.UUID(required=True)
        amount = graphene.Float(required=True)

    wallet = graphene.Field(lambda: WalletType)

    @login_required
    def mutate(root, info, wallet_id, amount):
        user = info.context.user
        wallet = Wallet.objects.get(wallet_id=wallet_id, user=user)
        wallet.topup(amount)
        Transaction.objects.create(
            from_wallet=None,
            to_wallet=wallet,
            amount=amount,
            transaction_type='topup',
            status='completed',
        )
        return TopUpWallet(wallet=wallet)


class SpendWallet(graphene.Mutation):
    class Arguments:
        wallet_id = graphene.UUID(required=True)
        amount = graphene.Float(required=True)

    wallet = graphene.Field(lambda: WalletType)

    @login_required
    def mutate(root, info, wallet_id, amount):
        user = info.context.user
        wallet = Wallet.objects.get(wallet_id=wallet_id, user=user)
        wallet.spend(amount)
        Transaction.objects.create(
            from_wallet=wallet,
            to_wallet=None,
            amount=amount,
            transaction_type='spend',
            status='completed',
        )
        return SpendWallet(wallet=wallet)


class CreateWallet(graphene.Mutation):
    class Arguments:
        wallet_type = graphene.String(required=True)
        user_email = graphene.String(required=True)  # new argument

    wallet = graphene.Field(lambda: WalletType)

    @login_required
    def mutate(root, info, wallet_type, user_email):
        # Only allow "pass" type
        if wallet_type != "pass":
            raise GraphQLError(
                "You can only create wallets of type 'pass' via API")

        from django.contrib.auth import get_user_model
        User = get_user_model()
        try:
            user = User.objects.get(email=user_email)
        except User.DoesNotExist:
            raise GraphQLError("User not found")

        # Get main wallet as parent
        try:
            parent_wallet = Wallet.objects.get(user=user, wallet_type="main")
        except Wallet.DoesNotExist:
            raise GraphQLError("Main wallet not found for this user")

        wallet = Wallet.objects.create(
            user=user,
            wallet_type=wallet_type,
            parent_wallet=parent_wallet
        )
        return CreateWallet(wallet=wallet)


class WalletMutation(graphene.ObjectType):
    topup_wallet = TopUpWallet.Field()
    spend_wallet = SpendWallet.Field()
    create_wallet = CreateWallet.Field()

import graphene
from wallet.models import Wallet, Transaction
from .queries import WalletType


class TopUpWallet(graphene.Mutation):
    class Arguments:
        wallet_id = graphene.UUID(required=True)
        amount = graphene.Float(required=True)

    wallet = graphene.Field(lambda: WalletType)

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


class WalletMutation(graphene.ObjectType):
    topup_wallet = TopUpWallet.Field()
    spend_wallet = SpendWallet.Field()

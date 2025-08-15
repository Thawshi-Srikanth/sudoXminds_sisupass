import graphene
from graphene_django import DjangoObjectType
from wallet.models import Wallet, Transaction

class WalletType(DjangoObjectType):
    class Meta:
        model = Wallet

class TransactionType(DjangoObjectType):
    class Meta:
        model = Transaction

class WalletQuery(graphene.ObjectType):
    wallets = graphene.List(WalletType)
    transactions = graphene.List(TransactionType)

    def resolve_wallets(root, info):
        user = info.context.user
        if user.is_anonymous:
            return Wallet.objects.none()
        return Wallet.objects.filter(user=user)

    def resolve_transactions(root, info):
        user = info.context.user
        if user.is_anonymous:
            return Transaction.objects.none()
        return Transaction.objects.filter(from_wallet__user=user) | Transaction.objects.filter(to_wallet__user=user)

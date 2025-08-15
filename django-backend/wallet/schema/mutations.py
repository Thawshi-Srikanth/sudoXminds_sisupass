# wallet/schema/mutations.py
import graphene
from graphql_jwt.decorators import login_required
from .types import WalletType
from ..models import Wallet

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


class WalletMutation(graphene.ObjectType):
    add_funds = AddFunds.Field()
    spend_funds = SpendFunds.Field()

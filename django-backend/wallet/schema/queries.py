from decimal import Decimal
import graphene
from graphene_django import DjangoObjectType
from wallet.models import Wallet, Transaction
from wallet.utils import login_required
from django.contrib.auth import get_user_model

User = get_user_model()

class UserType(DjangoObjectType):
    class Meta:
        model = User
        fields = ("id", "email") 

class WalletType(DjangoObjectType):
    balance = graphene.Float()
    class Meta:
        model = Wallet
        fields = ("wallet_id", "wallet_type", "balance", "parent_wallet", "user") 

    user = graphene.Field(UserType)  
    parent_wallet = graphene.Field(lambda: WalletType)  

    def resolve_balance(self, info):
        # Convert Decimal to float
        if isinstance(self.balance, Decimal):
            return float(self.balance)
        return self.balance

class TransactionType(DjangoObjectType):
    class Meta:
        model = Transaction

class WalletQuery(graphene.ObjectType):
    wallets = graphene.List(WalletType)
    transactions = graphene.List(TransactionType)

    @login_required
    def resolve_wallets(root, info):
        user = info.context.user
        if user.is_anonymous:
            return Wallet.objects.none()
        return Wallet.objects.filter(user=user)

    @login_required
    def resolve_transactions(root, info):
        user = info.context.user
        if user.is_anonymous:
            return Transaction.objects.none()
        return Transaction.objects.filter(from_wallet__user=user) | Transaction.objects.filter(to_wallet__user=user)

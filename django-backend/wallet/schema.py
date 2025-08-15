import graphene
from graphene_django import DjangoObjectType
from graphql_jwt.decorators import login_required
from .models import Wallet, Transaction


class WalletType(DjangoObjectType):
    class Meta:
        model = Wallet
        fields = ("wallet_id", "balance", "wallet_type",
                  "created_at", "updated_at")


class TransactionType(DjangoObjectType):
    class Meta:
        model = Transaction
        fields = ("transaction_id", "from_wallet", "to_wallet", "amount",
                  "transaction_type", "status", "transaction_date", "description")

# ------------------ QUERIES ------------------


class WalletQuery(graphene.ObjectType):
    my_wallets = graphene.List(WalletType)
    my_transactions = graphene.List(TransactionType)

    @login_required
    def resolve_my_wallets(root, info):
        user = info.context.user
        return Wallet.objects.filter(user=user)

    @login_required
    def resolve_my_transactions(root, info):
        user = info.context.user
        # Include transactions where the user's wallet is either from_wallet or to_wallet
        return Transaction.objects.filter(
            from_wallet__user=user
        ) | Transaction.objects.filter(
            to_wallet__user=user
        )

# ------------------ MUTATIONS ------------------


class AddFunds(graphene.Mutation):
    class Arguments:
        wallet_id = graphene.UUID(required=True)
        amount = graphene.Decimal(required=True)

    wallet = graphene.Field(WalletType)

    @login_required
    def mutate(self, info, wallet_id, amount):
        user = info.context.user
        try:
            wallet = Wallet.objects.get(wallet_id=wallet_id, user=user)
        except Wallet.DoesNotExist:
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
        try:
            wallet = Wallet.objects.get(wallet_id=wallet_id, user=user)
        except Wallet.DoesNotExist:
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

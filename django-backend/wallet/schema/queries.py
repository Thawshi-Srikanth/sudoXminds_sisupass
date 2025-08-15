# wallet/schema/queries.py
import graphene
from graphql_jwt.decorators import login_required
from .types import WalletType, TransactionType
from ..models import Wallet, Transaction

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
        return Transaction.objects.filter(from_wallet__user=user) | Transaction.objects.filter(to_wallet__user=user)

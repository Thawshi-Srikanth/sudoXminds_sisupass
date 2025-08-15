# wallet/schema/types.py
import graphene
from graphene_django import DjangoObjectType
from ..models import Wallet, Transaction

class WalletType(DjangoObjectType):
    class Meta:
        model = Wallet
        fields = ("wallet_id", "balance", "wallet_type", "created_at", "updated_at")


class TransactionType(DjangoObjectType):
    class Meta:
        model = Transaction
        fields = ("transaction_id", "from_wallet", "to_wallet", "amount",
                  "transaction_type", "status", "transaction_date", "description")

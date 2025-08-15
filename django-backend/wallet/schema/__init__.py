# wallet/schema/__init__.py
import graphene
from .queries import WalletQuery, PassQuery, TransactionQuery,
from .mutations import WalletMutation, PassMutation, PassDetailMutation


class Query(WalletQuery, PassQuery, TransactionQuery, graphene.ObjectType):
    pass


class Mutation(WalletMutation, PassMutation, PassDetailMutation, graphene.ObjectType):
    pass


schema = graphene.Schema(query=Query, mutation=Mutation)

# wallet/schema/__init__.py
import graphene
from .queries import WalletQuery, PassQuery
from .mutations import WalletMutation,PassMutation


class Query(WalletQuery, PassQuery,  graphene.ObjectType):
    pass


class Mutation(WalletMutation,PassMutation, graphene.ObjectType):
    pass


schema = graphene.Schema(query=Query, mutation=Mutation)

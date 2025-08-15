# wallet/schema/__init__.py
import graphene
from .queries import WalletQuery
from .mutations import WalletMutation


class Query(WalletQuery, graphene.ObjectType):
    pass


class Mutation(WalletMutation, graphene.ObjectType):
    pass


schema = graphene.Schema(query=Query, mutation=Mutation)

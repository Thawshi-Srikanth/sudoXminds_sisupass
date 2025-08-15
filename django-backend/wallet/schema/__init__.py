import graphene
from wallet.schema.queries import WalletQuery
from wallet.schema.mutations import WalletMutation


class Query(WalletQuery, graphene.ObjectType):
    pass


class Mutation(WalletMutation, graphene.ObjectType):
    pass


schema = graphene.Schema(query=Query, mutation=Mutation)

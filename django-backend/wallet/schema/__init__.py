import graphene
from wallet.schema.queries import WalletQuery
from wallet.schema.mutations import TokenAuthSocial, WalletMutation
import graphql_jwt 

class Query(WalletQuery, graphene.ObjectType):
    pass




class Mutation(WalletMutation, graphene.ObjectType):
    token_auth = graphql_jwt.ObtainJSONWebToken.Field()
    verify_token = graphql_jwt.Verify.Field()
    refresh_token = graphql_jwt.Refresh.Field()
    token_auth_social = TokenAuthSocial.Field() 


schema = graphene.Schema(query=Query, mutation=Mutation)

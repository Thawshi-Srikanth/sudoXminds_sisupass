import graphene
import graphql_jwt
from wallet.schema import WalletQuery, WalletMutation, PassQuery, PassMutation
from django.contrib.auth import get_user_model
from graphene_django import DjangoObjectType

User = get_user_model()


class UserType(DjangoObjectType):
    class Meta:
        model = User
        fields = ("id", "email", "username")


class AuthMutations(graphene.ObjectType):
    token_auth = graphql_jwt.ObtainJSONWebToken.Field()
    verify_token = graphql_jwt.Verify.Field()
    refresh_token = graphql_jwt.Refresh.Field()


class Query(WalletQuery, PassQuery, graphene.ObjectType):
    me = graphene.Field(lambda: UserType)

    def resolve_me(root, info):
        user = info.context.user
        if user.is_anonymous:
            raise Exception("Authentication Failure!")
        return user


class Mutation(WalletMutation, PassMutation, AuthMutations, graphene.ObjectType):
    pass


schema = graphene.Schema(query=Query, mutation=Mutation)

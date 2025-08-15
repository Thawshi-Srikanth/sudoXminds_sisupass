# auth/mutations.py
import graphene
from django.contrib.auth import authenticate
from rest_framework_simplejwt.tokens import RefreshToken
from django.contrib.auth import get_user_model

User = get_user_model()


class ObtainToken(graphene.Mutation):
    access = graphene.String()
    refresh = graphene.String()
    
    class Arguments:
        email = graphene.String(required=True)
        password = graphene.String(required=True)

    def mutate(self, info, email, password):
        user = authenticate(username=email, password=password)
        if not user:
            raise Exception("Invalid credentials")

        refresh = RefreshToken.for_user(user)
        return ObtainToken(
            access=str(refresh.access_token),
            refresh=str(refresh)
        )

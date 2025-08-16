from allauth.socialaccount.providers.google.views import GoogleOAuth2Adapter
from allauth.socialaccount.providers.oauth2.client import OAuth2Client
from dj_rest_auth.registration.views import SocialLoginView
from .models import CustomUser
from rest_framework.generics import CreateAPIView
from .serializers import RegisterSerializer
from rest_framework import status
from rest_framework_simplejwt.tokens import RefreshToken
from rest_framework.response import Response


class GoogleLogin(SocialLoginView):
    adapter_class = GoogleOAuth2Adapter


class RegisterView(CreateAPIView):
    queryset = CustomUser.objects.all()
    permission_classes = ()
    serializer_class = RegisterSerializer

    def post(self, request, *args, **kwargs):
        serializer = self.get_serializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        user = serializer.save()

        # Generate JWT tokens
        refresh = RefreshToken.for_user(user)
        access_token = str(refresh.access_token)

        response = Response(
            {
                "access": access_token,
                "refresh": str(refresh),
                "message": "User registered successfully",
            },
            status=status.HTTP_201_CREATED
        )

        # Set cookies
        response.set_cookie(
            key="jwt-auth", value=access_token, httponly=True, secure=False, samesite="Lax"
        )
        response.set_cookie(
            key="refresh-token", value=str(refresh), httponly=True, secure=False, samesite="Lax"
        )

        return response

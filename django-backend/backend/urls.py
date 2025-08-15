from django.contrib import admin
from django.urls import path, include
from rest_framework_simplejwt.views import (
    TokenObtainPairView,
    TokenRefreshView,
)
from graphene_django.views import GraphQLView
from django.views.decorators.csrf import csrf_exempt
from authentication.views import GoogleLogin, RegisterView
from graphql_jwt.decorators import jwt_cookie

urlpatterns = [
    path('admin/', admin.site.urls),

    path('api/auth/', include('rest_framework.urls')),
    path('api/token/', TokenObtainPairView.as_view(), name='token_obtain_pair'),
    path('api/token/refresh/', TokenRefreshView.as_view(), name='token_refresh'),
    path('accounts/', include('allauth.urls')),  # Google login

    path('auth/', include('dj_rest_auth.urls')),
    path('auth/registration/', RegisterView.as_view(), name="custom-register"),
    path('auth/google/', GoogleLogin.as_view(), name='google_login'),

    path(
        "graphql/",
        csrf_exempt(
            jwt_cookie(GraphQLView.as_view(graphiql=True))
        ),
        name="graphql"
    ),
]

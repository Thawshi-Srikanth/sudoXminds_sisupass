from allauth.socialaccount.helpers import complete_social_login
from allauth.socialaccount.providers.google.views import GoogleOAuth2Adapter
from allauth.socialaccount.models import SocialAccount
from rest_framework import status
from dj_rest_auth.registration.views import SocialLoginView

class GoogleLogin(SocialLoginView):
    adapter_class = GoogleOAuth2Adapter

    def post(self, request, *args, **kwargs):
        # This will handle Google login, user creation, and storing extra data in SocialAccount.extra_data
        response = super().post(request, *args, **kwargs)
        if response.status_code == status.HTTP_200_OK:
            user = self.request.user
            social_account = SocialAccount.objects.filter(user=user, provider='google').first()
            extra_data = social_account.extra_data if social_account else {}
            response.data['google_extra'] = extra_data
        return response
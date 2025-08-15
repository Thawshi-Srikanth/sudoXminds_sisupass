from rest_framework import serializers
from django.contrib.auth import get_user_model
User = get_user_model()


class GoogleAuthResponseSerializer(serializers.Serializer):
    sub = serializers.CharField()
    email = serializers.EmailField()
    email_verified = serializers.BooleanField()
    name = serializers.CharField()
    picture = serializers.URLField()
    given_name = serializers.CharField()
    family_name = serializers.CharField(required=False)


User = get_user_model()


# --- Google OAuth Response ---
class GoogleAuthResponseSerializer(serializers.Serializer):
    sub = serializers.CharField()  # Google user ID
    email = serializers.EmailField()
    email_verified = serializers.BooleanField()
    given_name = serializers.CharField(required=False, allow_blank=True)
    family_name = serializers.CharField(required=False, allow_blank=True)
    picture = serializers.URLField(required=False, allow_blank=True)


# --- User Serializer ---
class UserSerializer(serializers.ModelSerializer):
    token = serializers.SerializerMethodField(read_only=True)

    class Meta:
        model = User
        fields = [
            "id",
            "username",
            "email",
            "first_name",
            "last_name",
            "google_auth_enabled",
            "extras",
            "token",  # optional, can be returned for frontend convenience
        ]

    def get_token(self, obj):
        # Only return if context has request and token is generated
        return getattr(obj, "token", None)

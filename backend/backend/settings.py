from datetime import timedelta
from pathlib import Path
import dj_database_url
from decouple import config, Csv

from .unfold_settings import UNFOLD  # noqa: F401

# ---------------------
# Base paths
# ---------------------
BASE_DIR = Path(__file__).resolve().parent.parent

# ---------------------
# Security & Debug
# ---------------------
SECRET_KEY = config("SECRET_KEY")
DEBUG = config("DEBUG", default=False, cast=bool)
ALLOWED_HOSTS = config("ALLOWED_HOSTS", default="localhost,127.0.0.1", cast=Csv())

# ---------------------
# Database
# ---------------------
DB_STRING = config("DB_STRING")
DATABASES = {
    "default": dj_database_url.parse(DB_STRING, conn_max_age=600),
}

# ---------------------
# Installed apps
# ---------------------
INSTALLED_APPS = [
    "unfold",
    "unfold.contrib.filters",
    "unfold.contrib.forms",
    "unfold.contrib.inlines",
    "unfold.contrib.import_export",
    "unfold.contrib.guardian",
    "unfold.contrib.simple_history",
    "unfold.contrib.location_field",
    "unfold.contrib.constance",

    "django.contrib.admin",
    "django.contrib.auth",
    "django.contrib.contenttypes",
    "django.contrib.sessions",
    "django.contrib.messages",
    "django.contrib.staticfiles",

    "authentication",
    "wallet",
    "slots",

    "corsheaders",
    "graphql_jwt",

    "rest_framework",
    "rest_framework.authkey",
    "dj_rest_auth",

    "graphene_django",

    "django_celery_beat",
    "django_celery_results",

    "django.contrib.sites",
    "allauth",
    "allauth.account",
    "allauth.socialaccount",
    "dj_rest_auth.registration",

    "allauth.socialaccount.providers.google",
    "graphql_jwt.refresh_token.apps.RefreshTokenConfig",
]

# ---------------------
# Middleware
# ---------------------
MIDDLEWARE = [
    "corsheaders.middleware.CorsMiddleware",
    "django.middleware.security.SecurityMiddleware",
    "django.contrib.sessions.middleware.SessionMiddleware",
    "django.middleware.common.CommonMiddleware",
    "django.middleware.csrf.CsrfViewMiddleware",
    "django.contrib.auth.middleware.AuthenticationMiddleware",
    "django.contrib.messages.middleware.MessageMiddleware",
    "django.middleware.clickjacking.XFrameOptionsMiddleware",
    "allauth.account.middleware.AccountMiddleware",
]

# ---------------------
# URLs and WSGI
# ---------------------
ROOT_URLCONF = "backend.urls"
WSGI_APPLICATION = "backend.wsgi.application"

# ---------------------
# Auth
# ---------------------
AUTHENTICATION_BACKENDS = (
    "graphql_jwt.backends.JSONWebTokenBackend",
    "allauth.account.auth_backends.AuthenticationBackend",
    "django.contrib.auth.backends.ModelBackend",
)
AUTH_USER_MODEL = "authentication.CustomUser"
ACCOUNT_SIGNUP_FIELDS = ["email*", "password1*", "password2*"]
ACCOUNT_LOGIN_METHODS = {"phone", "email"}

# ---------------------
# JWT
# ---------------------
SIMPLE_JWT = {
    "ACCESS_TOKEN_LIFETIME": timedelta(
        minutes=config("JWT_ACCESS_LIFETIME_MIN", default=60, cast=int)
    ),
    "REFRESH_TOKEN_LIFETIME": timedelta(
        days=config("JWT_REFRESH_LIFETIME_DAYS", default=7, cast=int)
    ),
    "AUTH_HEADER_TYPES": ("Bearer",),
}

JWT_AUTH_COOKIE_SAMESITE = config("JWT_AUTH_COOKIE_SAMESITE", default="None")
JWT_AUTH_COOKIE_SECURE = config("JWT_AUTH_COOKIE_SECURE", default=True, cast=bool)

GRAPHQL_JWT = {
    "JWT_VERIFY_EXPIRATION": True,
    "JWT_LONG_RUNNING_REFRESH_TOKEN": True,
}

REST_AUTH = {
    "USE_JWT": True,
    "JWT_AUTH_COOKIE": "jwt-auth",
    "JWT_AUTH_REFRESH_COOKIE": "refresh-token",
    "JWT_AUTH_HTTPONLY": False,
}

# ---------------------
# Social login
# ---------------------
SOCIALACCOUNT_PROVIDERS = {
    "google": {
        "SCOPE": ["profile", "email"],
        "AUTH_PARAMS": {"access_type": "online"},
        "VERIFIED_EMAIL": True,
        "CLIENT_ID": config("GOOGLE_CLIENT_ID", default=""),
        "SECRET": config("GOOGLE_CLIENT_SECRET", default=""),
    }
}

# ---------------------
# CORS & CSRF
# ---------------------
CORS_ALLOW_HEADERS = [
    "accept", "accept-encoding", "authorization", "content-type",
    "dnt", "origin", "user-agent", "x-csrftoken", "x-requested-with",
]
CORS_ALLOW_CREDENTIALS = True

if DEBUG:
    CORS_ALLOW_ALL_ORIGINS = True
else:
    CORS_ALLOWED_ORIGINS = config(
        "CORS_ALLOWED_ORIGINS",
        default="http://localhost:3000,http://127.0.0.1:3000",
        cast=Csv()
    )

CSRF_TRUSTED_ORIGINS = config(
    "CSRF_TRUSTED_ORIGINS",
    default="http://localhost:3000",
    cast=Csv()
)

# ---------------------
# REST Framework
# ---------------------
REST_FRAMEWORK = {
    "DEFAULT_AUTHENTICATION_CLASSES": (
        "dj_rest_auth.jwt_auth.JWTCookieAuthentication",
    ),
    "DEFAULT_PERMISSION_CLASSES": [
        "rest_framework.permissions.DjangoModelPermissionsOrAnonReadOnly"
    ],
}

# ---------------------
# Templates
# ---------------------
TEMPLATES = [
    {
        "BACKEND": "django.template.backends.django.DjangoTemplates",
        "DIRS": [BASE_DIR / "templates"],
        "APP_DIRS": True,
        "OPTIONS": {
            "context_processors": [
                "django.template.context_processors.request",
                "django.contrib.auth.context_processors.auth",
                "django.contrib.messages.context_processors.messages",
            ],
        },
    },
]

# ---------------------
# Password validation
# ---------------------
AUTH_PASSWORD_VALIDATORS = [
    {"NAME": "django.contrib.auth.password_validation.UserAttributeSimilarityValidator"},
    {"NAME": "django.contrib.auth.password_validation.MinimumLengthValidator"},
    {"NAME": "django.contrib.auth.password_validation.CommonPasswordValidator"},
    {"NAME": "django.contrib.auth.password_validation.NumericPasswordValidator"},
]

# ---------------------
# Internationalization
# ---------------------
LANGUAGE_CODE = "en-us"
TIME_ZONE = "UTC"
USE_I18N = True
USE_TZ = True

# ---------------------
# Static files
# ---------------------
STATIC_URL = config("STATIC_URL", default="static/")
STATICFILES_DIRS = [BASE_DIR / config("STATICFILES_DIR", default="static")]

DEFAULT_AUTO_FIELD = "django.db.models.BigAutoField"

# ---------------------
# Email
# ---------------------
EMAIL_BACKEND = (
    "django.core.mail.backends.console.EmailBackend" if DEBUG else "django.core.mail.backends.smtp.EmailBackend"
)
EMAIL_HOST = config("EMAIL_HOST", default="localhost")
EMAIL_PORT = config("EMAIL_PORT", default=587, cast=int)
EMAIL_HOST_USER = config("EMAIL_HOST_USER", default="")
EMAIL_HOST_PASSWORD = config("EMAIL_HOST_PASSWORD", default="")
EMAIL_USE_TLS = config("EMAIL_USE_TLS", default=True, cast=bool)
DEFAULT_FROM_EMAIL = config("DEFAULT_FROM_EMAIL", default="webmaster@localhost")

# ---------------------
# Celery
# ---------------------
CELERY_BROKER_URL = config("CELERY_BROKER_URL", default="redis://localhost:6379/0")
CELERY_RESULT_BACKEND = config("CELERY_RESULT_BACKEND", default="redis://localhost:6379/0")
CELERY_BEAT_SCHEDULER = "django_celery_beat.schedulers:DatabaseScheduler"

# ---------------------
# Site
# ---------------------
SITE_ID = config("SITE_ID", default=3, cast=int)

# ---------------------
# Redirect URLs
# ---------------------
LOGIN_REDIRECT_URL = "/"
LOGOUT_REDIRECT_URL = "/"

# ---------------------
# Graphene
# ---------------------
GRAPHENE = {
    "SCHEMA": "backend.schema.schema",
    "MIDDLEWARE": ["graphql_jwt.middleware.JSONWebTokenMiddleware"],
}

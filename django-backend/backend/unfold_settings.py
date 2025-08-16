from django.templatetags.static import static
from django.urls import reverse_lazy
from django.utils.translation import gettext_lazy as _

UNFOLD = {
    "SITE_TITLE": "SiSu Pass",
    "SITE_HEADER": "SiSu Pass Admin",
    "SITE_SUBHEADER": "Identity · Wallet · Smart Slots",
    "SITE_ICON": {
        "light": lambda request: static("icons/sisu-light.svg"),
        "dark": lambda request: static("icons/sisu-dark.svg"),
    },
    "SITE_LOGO": {
        "light": lambda request: static("logos/sisu-logo-light.svg"),
        "dark": lambda request: static("logos/sisu-logo-dark.svg"),
    },
    "SITE_URL": "/",
    "THEME": 'light',
    "COMMAND": {
        "search_models": True,  
        "search_callback": "utils.search_callback",
        "show_history": True, 
    },
    "COLORS": {
        "base": {
            "50": "249, 250, 251",
            "100": "243, 244, 246",
            "200": "229, 231, 235",
            "300": "209, 213, 219",
            "400": "156, 163, 175",
            "500": "107, 114, 128",
            "600": "75, 85, 99",
            "700": "55, 65, 81",
            "800": "31, 41, 55",
            "900": "17, 24, 39",
            "950": "3, 7, 18",
        },
        "primary": {
            "50":  "236, 252, 245",
            "100": "209, 250, 229",
            "200": "167, 243, 208",
            "300": "110, 231, 183",
            "400": "52, 211, 153",
            "500": "27, 196, 125",
            "600": "22, 163, 104",
            "700": "18, 130, 83",
            "800": "14, 97, 62",
            "900": "11, 65, 41",
        },
        "font": {
            "subtle-light": "var(--color-base-500)",
            "subtle-dark": "var(--color-base-400)",
            "default-light": "var(--color-base-600)",
            "default-dark": "var(--color-base-300)",
            "important-light": "var(--color-base-900)",
            "important-dark": "var(--color-base-100)",
        },
    },

    # Optional tweaks
    "BORDER_RADIUS": "8px",
    "STYLES": [],   # no external css
    "SCRIPTS": [],  # no external js
    "SIDEBAR": {
        "navigation": [
            {
                "title": _("Users & Auth"),
                "icon": "people",
                "items": [
                    {"title": _("Users"), "link": reverse_lazy("admin:authentication_customuser_changelist")},
                ],
            },
            {
                "title": _("Wallets & Passes"),
                "icon": "account_balance_wallet",
                "items": [
                    {"title": _("Wallets"), "link": reverse_lazy("admin:wallet_wallet_changelist")},
                    {"title": _("Transactions"), "link": reverse_lazy("admin:wallet_transaction_changelist")},
                    {"title": _("Pass Details"), "link": reverse_lazy("admin:wallet_passdetails_changelist")},
                ],
            },
            {
                "title": _("Slots & Bookings"),
                "icon": "event",
                "items": [
                    {"title": _("Slot Types"), "link": reverse_lazy("admin:slots_slottype_changelist")},
                    {"title": _("Slots"), "link": reverse_lazy("admin:slots_slot_changelist")},
                    {"title": _("Bookings"), "link": reverse_lazy("admin:slots_booking_changelist")},
                ],
            },
        ],
    },
}


# === OPTIONAL CALLBACKS ===

def dashboard_callback(request, context):
    from wallet.models import Transaction, Wallet
    context.update({
        "total_wallets": Wallet.objects.count(),
        "transactions_today": Transaction.objects.filter(
            transaction_date__date=request.today()
        ).count(),
    })
    return context


def environment_callback(request):
    return ["Production", "success"]


def badge_callback(request):
    from slots.models import Booking
    return Booking.objects.filter(status="pending").count()


def permission_callback(request):
    return request.user.is_superuser

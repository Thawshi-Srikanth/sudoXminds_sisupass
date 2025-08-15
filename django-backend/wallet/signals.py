from django.db.models.signals import post_save
from django.dispatch import receiver
from django.conf import settings
from wallet.models import Wallet

@receiver(post_save, sender=settings.AUTH_USER_MODEL)
def create_default_wallet(sender, instance, created, **kwargs):
    if created:
        # Create default wallet for new users
        Wallet.objects.create(user=instance, wallet_type="main")

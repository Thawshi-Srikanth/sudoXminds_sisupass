import uuid
from django.db import models
from django.conf import settings
from django.core.exceptions import ValidationError


class Wallet(models.Model):

    WALLET_TYPES = [
        ('main', 'Main'),
        ('pass', 'Pass'),
        ('vendor', 'Vendor'),
    ]

    wallet_id = models.UUIDField(
        primary_key=True, default=uuid.uuid4, editable=False)
    user = models.ForeignKey(settings.AUTH_USER_MODEL,
                             on_delete=models.CASCADE, related_name="wallets")
    balance = models.DecimalField(default=0, max_digits=20, decimal_places=2)
    wallet_type = models.CharField(max_length=20, choices=WALLET_TYPES)

    parent_wallet = models.ForeignKey(
        'self', on_delete=models.SET_NULL, null=True, blank=True)
    access_qr_code = models.CharField(
        max_length=255, unique=True, null=True, blank=True)
    access_nfc_code = models.CharField(
        max_length=255, unique=True, null=True, blank=True)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    def __str__(self):
        return f"{self.user.email} - {self.wallet_type} ({self.balance})"

    def topup(self, amount):
        if amount <= 0:
            raise ValidationError("Amount must be positive")
        self.balance += amount
        self.save()
        return self.balance

    def spend(self, amount):
        if amount <= 0:
            raise ValidationError("Amount must be positive")
        if amount > self.balance:
            raise ValidationError("Insufficient balance")
        self.balance -= amount
        self.save()
        return self.balance

    class Meta:
        constraints = [
            models.UniqueConstraint(
                fields=['user', 'wallet_type'],
                condition=models.Q(wallet_type='main'),
                name='unique_main_wallet_per_user'
            )
        ]


class Transaction(models.Model):
    TRANSACTION_TYPES = [
        ('topup', 'Top Up'),
        ('spend', 'Spend'),
        ('transfer', 'Transfer'),
    ]
    STATUS_CHOICES = [
        ('pending', 'Pending'),
        ('completed', 'Completed'),
        ('failed', 'Failed'),
    ]

    transaction_id = models.UUIDField(
        primary_key=True, default=uuid.uuid4, editable=False)
    from_wallet = models.ForeignKey(
        Wallet, related_name='outgoing_transactions', on_delete=models.SET_NULL, null=True, blank=True)
    to_wallet = models.ForeignKey(
        Wallet, related_name='incoming_transactions', on_delete=models.SET_NULL, null=True, blank=True)
    amount = models.DecimalField(max_digits=20, decimal_places=2)
    transaction_date = models.DateTimeField(auto_now_add=True)
    transaction_type = models.CharField(
        max_length=50, choices=TRANSACTION_TYPES)
    status = models.CharField(
        max_length=20, choices=STATUS_CHOICES, default='completed')
    description = models.TextField(null=True, blank=True)

    def __str__(self):
        return f"{self.transaction_type} - {self.amount}"

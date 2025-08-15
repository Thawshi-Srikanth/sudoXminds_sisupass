import random
from decimal import Decimal
from django.core.management.base import BaseCommand
from wallet.models import Wallet, Transaction

class Command(BaseCommand):
    help = "Populate sample transactions for existing wallets"

    def handle(self, *args, **options):
        # Step 1: Get all wallets
        wallets = Wallet.objects.all()
        if not wallets.exists():
            self.stdout.write(self.style.ERROR("No wallets found!"))
            return

        transaction_types = ['topup', 'spend', 'transfer']
        status_choices = ['completed', 'pending', 'failed']

        created_count = 0

        for wallet in wallets:
            # Decide how many transactions to create for this wallet
            num_tx = random.randint(1, 5)

            for _ in range(num_tx):
                # Choose transaction type
                tx_type = random.choice(transaction_types)
                
                # Skip topup if wallet is not main
                if tx_type == 'topup' and wallet.wallet_type != 'main':
                    continue

                amount = round(random.uniform(10, 500), 2)
                status = random.choice(status_choices)

                if tx_type == 'topup':
                    from_wallet = None
                    to_wallet = wallet
                elif tx_type == 'spend':
                    from_wallet = wallet
                    to_wallet = None
                else:  # transfer
                    from_wallet = wallet
                    # Pick a random target wallet different from this wallet
                    possible_targets = wallets.exclude(wallet_id=wallet.wallet_id)
                    if possible_targets.exists():
                        to_wallet = random.choice(possible_targets)
                    else:
                        to_wallet = None

                # Create transaction
                Transaction.objects.create(
                    from_wallet=from_wallet,
                    to_wallet=to_wallet,
                    amount=Decimal(amount),
                    transaction_type=tx_type,
                    status=status,
                    description=f"Sample {tx_type} transaction"
                )
                created_count += 1

        self.stdout.write(self.style.SUCCESS(f"Created {created_count} transactions successfully!"))

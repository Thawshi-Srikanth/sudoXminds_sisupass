from unfold.dataclasses import SearchResult
from wallet.models import Wallet, PassDetails, Transaction
from authentication.models import CustomUser


def search_callback(request, search_term: str):
    results = []

    # Clean up the search term
    search_term = search_term.strip()
    if not search_term or len(search_term) < 2:
        return results

    # --- Search Users by email or name ---
    users = CustomUser.objects.filter(email__icontains=search_term)[:5]
    for user in users:
        results.append(SearchResult(
            title=user.get_full_name() or user.email,
            description=f"User email: {user.email}",
            link=f"/admin/authentication/customuser/{user.id}/change/",
            icon="person",
        ))

    # --- Search Wallets by user email ---
    wallets = Wallet.objects.filter(user__email__icontains=search_term)[:5]
    for wallet in wallets:
        results.append(SearchResult(
            title=f"{wallet.user.email} - Wallet #{wallet.id}",
            description=f"Balance: {wallet.balance}",
            link=f"/admin/wallet/wallet/{wallet.id}/change/",
            icon="account_balance_wallet",
        ))

    # --- Search PassDetails by category ---
    passes = PassDetails.objects.filter(category__icontains=search_term)[:5]
    for p in passes:
        results.append(SearchResult(
            title=f"{p.category} - Wallet #{p.wallet.id}",
            description="PassDetails",
            link=f"/admin/wallet/passdetails/{p.id}/change/",
            icon="badge",
        ))

    # --- Search Transactions by type ---
    transactions = Transaction.objects.filter(transaction_type__icontains=search_term)[:5]
    for t in transactions:
        results.append(SearchResult(
            title=f"{t.transaction_type} - Wallet #{t.wallet.id}",
            description=f"Amount: {t.amount} | Status: {t.status}",
            link=f"/admin/wallet/transaction/{t.id}/change/",
            icon="swap_horiz",
        ))

    return results

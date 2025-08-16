from django.shortcuts import render
from django.contrib.admin.views.decorators import staff_member_required
from wallet.models import Wallet, Transaction, PassDetails
from django.db.models import Count
from django.utils.timezone import now, timedelta

@staff_member_required
def dashboard_view(request):
    today = now().date()
    seven_days_ago = today - timedelta(days=7)

    # Aggregate data
    transaction_stats = (
        Transaction.objects.filter(transaction_date__gte=seven_days_ago)
        .values('transaction_date__date')
        .annotate(count=Count('id'))
        .order_by('transaction_date__date')
    )
    wallet_stats = (
        Wallet.objects.filter(created_at__gte=seven_days_ago)
        .values('created_at__date')
        .annotate(count=Count('id'))
        .order_by('created_at__date')
    )
    pass_stats = (
        PassDetails.objects.filter(wallet__created_at__gte=seven_days_ago)
        .values('wallet__created_at__date')
        .annotate(count=Count('id'))
        .order_by('wallet__created_at__date')
    )

    context = {
        'page_title': 'Wallet Dashboard',
        'transactions': transaction_stats,
        'wallets': wallet_stats,
        'passes': pass_stats,
    }
    return render(request, 'admin/dashboard.html', context)

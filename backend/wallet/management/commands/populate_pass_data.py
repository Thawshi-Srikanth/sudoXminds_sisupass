import random
from datetime import date, timedelta
from django.core.management.base import BaseCommand
from django.contrib.auth import get_user_model
from wallet.models import Wallet, PassDetails, LocationType, Location, PassCategory

User = get_user_model()


class Command(BaseCommand):
    help = "Populate sample users, wallets, passes, locations, and categories"

    def handle(self, *args, **kwargs):
        NUM_USERS = 5  # Number of test users to create

        # -------------------------------
        # 0. Clear all existing test data
        # -------------------------------
        Wallet.objects.all().delete()
        PassDetails.objects.all().delete()
        Location.objects.all().delete()
        LocationType.objects.all().delete()
        PassCategory.objects.all().delete()
        User.objects.filter(email__endswith="@sisupass.com").delete()
        self.stdout.write(self.style.WARNING(
            "Cleared existing test users, wallets, passes, locations, and categories."))

        # -------------------------------
        # 1. Create Location Types
        # -------------------------------
        location_type_names = ["Bus Station",
                               "Train Station", "Restaurant", "Library"]
        location_types = {}
        for lt_name in location_type_names:
            lt, _ = LocationType.objects.get_or_create(name=lt_name)
            location_types[lt_name] = lt
            self.stdout.write(self.style.SUCCESS(
                f"Created LocationType: {lt_name}"))

        # -------------------------------
        # 2. Create 10 locations per type
        # -------------------------------
        locations = {}
        for lt_name, lt in location_types.items():
            locations[lt_name] = []
            for i in range(1, 11):
                loc_name = f"{lt_name} {i}"
                loc = Location.objects.create(
                    name=loc_name,
                    address=f"{loc_name} Address",
                    latitude=round(random.uniform(-10, 10), 6),
                    longitude=round(random.uniform(75, 85), 6),
                    location_type=lt
                )
                locations[lt_name].append(loc)
            self.stdout.write(self.style.SUCCESS(
                f"Created 10 locations for {lt_name}"))

        # -------------------------------
        # 3. Create Pass Categories
        # -------------------------------
        category_mapping = {
            "Bus Pass": ["Bus Station"],
            "Train Pass": ["Train Station"],
            "Food Pass": ["Restaurant"],
            "Library Pass": ["Library"],
        }
        categories = {}
        for cat_name, allowed_lt_names in category_mapping.items():
            category = PassCategory.objects.create(name=cat_name)
            allowed_lts = [location_types[lt_name]
                           for lt_name in allowed_lt_names]
            category.allowed_location_types.set(allowed_lts)
            categories[cat_name] = category
            self.stdout.write(self.style.SUCCESS(
                f"Created PassCategory: {cat_name}"))

        # -------------------------------
        # 4. Create multiple users with wallets and passes
        # -------------------------------
        for u in range(1, NUM_USERS + 1):
            username = f"user{u}"
            email = f"user{u}@sisupass.com"
            password = "password123"

            user, _ = User.objects.get_or_create(
                username=username,
                defaults={"email": email}
            )
            user.set_password(password)
            user.save()
            self.stdout.write(self.style.SUCCESS(f"Created user: {username}"))

            # -------------------------------
            # 5. Create main wallet (get_or_create to avoid duplicates)
            # -------------------------------
            main_wallet, created = Wallet.objects.get_or_create(
                user=user,
                wallet_type="main",
                defaults={"balance": random.randint(100, 1000)}
            )
            if not created:
                main_wallet.balance = random.randint(100, 1000)
                main_wallet.save()
            self.stdout.write(self.style.SUCCESS(
                f"Created main wallet for {username}"))

            # -------------------------------
            # 6. Create pass wallets and PassDetails
            # -------------------------------
            for cat_name, category in categories.items():
                pass_wallet = Wallet.objects.create(
                    user=user,
                    wallet_type="pass",
                    parent_wallet=main_wallet,
                    balance=0
                )

                if cat_name in ["Bus Pass", "Train Pass"]:
                    from_loc = locations[category_mapping[cat_name][0]][0]
                    to_loc = locations[category_mapping[cat_name][0]][1]
                    pass_obj = PassDetails.objects.create(
                        wallet=pass_wallet,
                        category=category,
                        from_location=from_loc,
                        to_location=to_loc,
                        start_date=date.today(),
                        end_date=date.today() + timedelta(days=365)
                    )
                else:  # Food or Library
                    allowed_locs = locations[category_mapping[cat_name][0]][:5]
                    pass_obj = PassDetails.objects.create(
                        wallet=pass_wallet,
                        category=category,
                        start_date=date.today(),
                        end_date=date.today() + timedelta(days=365)
                    )
                    pass_obj.allowed_locations.set(allowed_locs)

                self.stdout.write(self.style.SUCCESS(
                    f"Created {cat_name} for {username}"))

        self.stdout.write(self.style.SUCCESS(
            f"Successfully created {NUM_USERS} test users with wallets and passes!"))

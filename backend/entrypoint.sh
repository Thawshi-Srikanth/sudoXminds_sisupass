python manage.py migrate

echo "Creating superuser..."
echo "from django.contrib.auth import get_user_model; User = get_user_model(); \
      User.objects.filter(username='admin').exists() or \
      User.objects.create_superuser('admin','admin@sisupass.com','admin')" | python manage.py shell


echo "Running custom management commands..."
python manage.py populate_pass_data
python manage.py populate_transactions
python manage.py populate_slot_types
python manage.py populate_slots
python manage.py populate_schedules


python manage.py runserver 0.0.0.0:8000

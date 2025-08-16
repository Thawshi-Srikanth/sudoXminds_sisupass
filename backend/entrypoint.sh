#!/bin/bash

# Wait for Postgres
echo "Waiting for postgres..."
python - <<END
import time, psycopg2, os
while True:
    try:
        conn = psycopg2.connect(
            dbname=os.environ['POSTGRES_DB'],
            user=os.environ['POSTGRES_USER'],
            password=os.environ['POSTGRES_PASSWORD'],
            host='postgres'
        )
        conn.close()
        break
    except:
        time.sleep(0.5)
END
echo "Postgres started"

# Apply migrations
python manage.py migrate

# Create superuser if it doesn't exist
echo "Creating superuser..."
echo "from django.contrib.auth import get_user_model; User = get_user_model(); \
      User.objects.filter(username='admin').exists() or \
      User.objects.create_superuser('admin','admin@sisupass.com','admin')" | python manage.py shell

# Run custom management commands to populate database
echo "Running custom management commands..."
python manage.py populate_pass_data
python manage.py populate_transactions
python manage.py populate_slot_types
python manage.py populate_slots
python manage.py populate_schedules



# Collect static files (optional)
python manage.py collectstatic --noinput

# Start Django server with gunicorn
gunicorn backend.wsgi:application --bind 0.0.0.0:8000

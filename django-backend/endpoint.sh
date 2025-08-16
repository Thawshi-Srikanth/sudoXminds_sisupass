#!/bin/bash

# Wait for Postgres
echo "Waiting for postgres..."
while ! nc -z postgres 5432; do
  sleep 0.1
done
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
python manage.py populate_initial_data  # <-- replace with your actual command(s)

# Collect static files (optional)
python manage.py collectstatic --noinput

# Start Django server with gunicorn
gunicorn backend.wsgi:application --bind 0.0.0.0:8000

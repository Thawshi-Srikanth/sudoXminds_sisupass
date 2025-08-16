# SiSu Pass Developement Quide

## Prerequisites

- **Python 3.10+** ([Download Python](https://www.python.org/downloads/))
- **pip** (comes with Python)
- **Virtualenv** (optional but recommended)
- **Database** (PostgreSQL/MySQL/SQLite)
- **Redis** (for Celery task queue)
- **Git** (optional, for version control)

Install Redis:

- **Windows:** Use [Redis for Windows](https://github.com/MicrosoftArchive/redis)
- **Linux/macOS:**

```bash
sudo apt install redis-server
brew install redis
```

Start Redis server:

```bash
redis-server
```

---

## 1. Clone the Project

```bash
git clone https://github.com/<username>/sudoXminds_sisupass
cd sudoXminds_sisupass
```

## 2. Set Up Virtual Environment

```bash
python -m venv venv
```

Activate it:

```bash
source venv/bin/activate
```

```

python -m venv venv

```

Activate it:

- **Windows:**

```bash
venv\Scripts\activate
```

- **Linux/macOS:**

```bash
source venv/bin/activate
```

## 3. Install Dependencies

```bash
pip install -r requirements.txt
```

## 4. Configure Environment Variables

Copy the contents of `.env.template` to a new file named `.env`:

```
SECRET_KEY=your_secret_key
DEBUG=True
DATABASE_URL=postgres://user:password@localhost:5432/dbname
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_HOST_USER=iam@gmail.com
EMAIL_HOST_PASSWORD=******************
EMAIL_USE_TLS=True
DEFAULT_FROM_EMAIL="SiSu Pass <iam@gmail.com>"
```

_(Update according to your `.env.template` or project settings)_

## 5. Apply Database Migrations

```bash
python manage.py migrate
```

## 6. Create Superuser (Admin Account)

```bash
python manage.py createsuperuser
```

## 7. Run the Development Server

```bash
python manage.py runserver
```

Open in browser:

```
http://127.0.0.1:8000/
```

## 8. Other Steps

- **Collect static files** (for production):

```bash
python manage.py collectstatic
```

- **Run Celery worker**:

Note: Celery does have known setup and usage challenges on Windows.

```bash
celery -A backend worker --loglevel=info
```

- **Run Celery beat scheduler**:

```bash
celery -A backend beat --loglevel=info
```

FROM mcr.microsoft.com/vscode/devcontainers/base:ubuntu

# Install system dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    golang-go postgresql-client curl git build-essential python3 python3-pip python3-venv libpq-dev python3-dev \
  && curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
  && apt-get install -y nodejs \
  && npm install -g pnpm \
  && npm cache clean --force \
  && rm -rf /var/lib/apt/lists/*


WORKDIR /workspace

EXPOSE 3000 5432 8080

CMD ["sleep", "infinity"]

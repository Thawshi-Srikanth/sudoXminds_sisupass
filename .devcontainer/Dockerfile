FROM mcr.microsoft.com/vscode/devcontainers/base:ubuntu

RUN apt-get update && apt-get install -y --no-install-recommends \
    golang-go postgresql-client curl git build-essential \
  && curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
  && apt-get install -y nodejs \
  && npm install -g pnpm \
  && npm cache clean --force \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /workspace

EXPOSE 3000 5432 8080

CMD ["sleep", "infinity"]

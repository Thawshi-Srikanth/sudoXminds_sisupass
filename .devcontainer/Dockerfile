FROM mcr.microsoft.com/vscode/devcontainers/go:1.24-bullseye AS build

RUN apt-get update && apt-get install -y --no-install-recommends postgresql-client

FROM mcr.microsoft.com/vscode/devcontainers/go:1.24-bullseye

COPY --from=build /usr/bin/psql /usr/bin/psql

WORKDIR /workspace

EXPOSE 3000 5432 8080
CMD [ "/bin/bash" ]
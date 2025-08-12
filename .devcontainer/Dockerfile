# Base image
FROM mcr.microsoft.com/devcontainers/base:bookworm

# Install curl & other utils
RUN apt-get update && apt-get install -y curl git unzip && rm -rf /var/lib/apt/lists/*

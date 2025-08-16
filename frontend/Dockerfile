FROM node:20

WORKDIR /app

# Install pnpm globally
RUN npm install -g pnpm

# Copy package files first to leverage Docker cache
COPY package.json pnpm-lock.yaml* ./

# Install dependencies
RUN pnpm install

# Copy the rest of the code
COPY . .

EXPOSE 3000

# Use CMD that works well with bind mounts
CMD ["pnpm", "dev"]

FROM node:20

WORKDIR /app

# Copy package files first for dependency caching
COPY package.json pnpm-lock.yaml ./

# Install pnpm globally
RUN npm install -g pnpm

# Install dependencies
RUN pnpm install

# Copy the rest of the frontend code
COPY . .

EXPOSE 3000

# Start development server
CMD ["pnpm", "run", "dev"]

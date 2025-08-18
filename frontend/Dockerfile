# Use Node 20 for development
FROM node:20

# Set working directory
WORKDIR /app

# Install pnpm globally
RUN npm install -g pnpm

# Copy package files first for caching
COPY package.json pnpm-lock.yaml* ./

# Install all dependencies (dev + prod)
RUN pnpm install

# Copy the rest of your code
COPY . .

# Expose the dev server port
EXPOSE 3000

# Start development server
CMD ["pnpm", "dev"]

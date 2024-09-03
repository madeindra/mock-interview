# Build stage
FROM node:20-alpine AS builder

# Set env for backend url
ARG BACKEND_URL=http://0.0.0.0:8080
ENV VITE_BACKEND_URL=${BACKEND_URL}

# Install pnpm
RUN npm install -g pnpm

WORKDIR /app

# Copy pnpm-lock.yaml file
COPY pnpm-lock.yaml ./

# Install dependencies
RUN pnpm fetch

# Copy the rest of the code
COPY . .

# Install dependencies
RUN pnpm install --frozen-lockfile

# Build the app
RUN pnpm run build

# Final stage
FROM node:20-alpine

WORKDIR /app

# Copy built assets from builder stage
COPY --from=builder /app/dist ./dist

# Copy package.json and pnpm-lock.yaml
COPY package.json pnpm-lock.yaml ./

# Install production dependencies only
RUN npm install -g pnpm && pnpm install --prod

# Install a simple HTTP server for serving static content
RUN npm install -g serve

EXPOSE 3000

CMD ["serve", "-s", "dist", "-l", "3000"]

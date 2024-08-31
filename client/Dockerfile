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
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
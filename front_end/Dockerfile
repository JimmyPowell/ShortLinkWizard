# Stage 1: Build the Next.js frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy the rest of the frontend code
COPY . .

# Build the Next.js app
RUN npm run build

# Stage 2: Build the Go backend
FROM golang:1.20-alpine AS backend-builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the backend code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 3: Create the final image
FROM alpine:3.14

WORKDIR /app

# Copy the built Next.js app from the frontend-builder stage
COPY --from=frontend-builder /app/.next ./.next
COPY --from=frontend-builder /app/public ./public
COPY --from=frontend-builder /app/package.json ./package.json
COPY --from=frontend-builder /app/node_modules ./node_modules

# Copy the built Go binary from the backend-builder stage
COPY --from=backend-builder /app/main .

# Install Node.js
RUN apk add --no-cache nodejs npm

# Expose the port your app runs on
EXPOSE 3000

# Start both the Next.js app and the Go backend
CMD ["sh", "-c", "npm start & ./main"]


# Stage 1: Build Svelte frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend

COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build

# Stage 2: Build Go backend
FROM golang:1.25-alpine AS backend-builder
WORKDIR /app/backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go

# Stage 3: Final lightweight container image
FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy Go server binary
COPY --from=backend-builder /app/backend/server /app/server

# Copy Svelte production build assets
COPY --from=frontend-builder /app/frontend/dist /app/frontend/dist

# Expose port
EXPOSE 8080

# Configure environment variables
ENV PORT=8080
ENV STATIC_DIR=/app/frontend/dist

# Run Go server
CMD ["/app/server"]


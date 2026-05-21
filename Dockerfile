# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app
# Since Render is starting inside the frontend folder, package.json is at the immediate root
COPY package.json package-lock.json ./
RUN npm ci
COPY . ./
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.25-alpine AS backend-builder
WORKDIR /app
# Reach up to the sibling backend folder
COPY ../backend/go.mod ../backend/go.sum ./backend/
RUN cd backend && go mod download
COPY ../backend/ ./backend/
RUN cd backend && CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Stage 3: Final Runner Container
FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=backend-builder /app/backend/main .
COPY --from=frontend-builder /app/dist ./dist

EXPOSE 8080
CMD ["./main"]

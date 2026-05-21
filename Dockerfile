# Stage 1: Build Svelte Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /workspace
COPY frontend/package*.json ./frontend/
RUN cd frontend && npm ci
COPY frontend/ ./frontend/
RUN cd frontend && npm run build

# Stage 2: Build Go Backend
FROM golang:1.25-alpine AS backend-builder
WORKDIR /workspace
COPY backend/go.mod backend/go.sum ./backend/
RUN cd backend && go mod download
COPY backend/ ./backend/
RUN cd backend && CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Stage 3: Final Production Runner Container
FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=backend-builder /workspace/backend/main .
COPY --from=frontend-builder /workspace/frontend/dist ./dist

EXPOSE 8080
CMD ["./main"]

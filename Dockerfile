# Stage 1: Build Svelte Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Stage 2: Build Go Backend
FROM golang:1.25-alpine AS backend-builder
WORKDIR /go/src/app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Stage 3: Final Production Runner Container
FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=backend-builder /go/src/app/main .
COPY --from=frontend-builder /app/dist ./dist

EXPOSE 8080
CMD ["./main"]

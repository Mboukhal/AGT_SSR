

# ---- BUILD UI ----
FROM oven/bun:1.3-alpine AS ui
WORKDIR /app
ENV APP_ENV=production
ENV NODE_ENV=production

COPY ./ui .
RUN bun install
RUN bun run build

# ---- BUILD BACKEND ----
FROM golang:1.26-alpine AS builder
WORKDIR /app
ENV APP_ENV=production

COPY --exclude=./ui . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go mod tidy
RUN go build -ldflags="-s -w" -o ./dist/server ./cmd/main.go

# ---- FINAL (scratch) ----
FROM scratch
# FROM golang:1.26-alpine

WORKDIR /app
ENV APP_ENV=production

COPY --from=ui /app/build ./ui

COPY --from=builder /app/dist/server .

EXPOSE 3000

CMD ["./server"]




# ---- BUILD UI ----
FROM oven/bun:1.3-alpine AS ui
WORKDIR /app
ENV APP_ENV=production
ENV NODE_ENV=production

COPY ./ui .
RUN bun install
RUN bun run build


# ---- BUILD sqlc ----
FROM golang:1.26-alpine AS sqlc
WORKDIR /app
ENV APP_ENV=production

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN /go/bin/sqlc generate

# ---- BUILD BACKEND ----
FROM golang:1.26-alpine AS builder
WORKDIR /app
ENV APP_ENV=production
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

COPY --exclude=./ui . .

RUN go mod download


# COPY --from=templ /app/ui/src/ /app/ui/src/
COPY --from=sqlc /app/internal/db/ /app/internal/db/
# COPY --from=sqlc /go/bin/sqlc /go/bin/sqlc
# RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
# RUN /go/bin/sqlc generate


RUN go build -ldflags="-s -w" -o ./dist/server ./cmd/main.go

# ---- FINAL (scratch) ----
FROM scratch
# FROM golang:1.26-alpine

WORKDIR /app
ENV APP_ENV=production

COPY --from=ui /app/dist ./static

COPY --from=builder /app/dist/server .

EXPOSE 3000

CMD ["./server"]


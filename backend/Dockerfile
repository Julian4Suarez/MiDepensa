# ── Stage 1: Build ──────────────────────────────────────────
FROM docker.io/library/golang:1.26-alpine AS builder

WORKDIR /app

# Dependency files first so this layer stays cached while only code changes.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# -s -w strip the symbol table and DWARF info (~30% smaller binary).
ARG APP_COMMIT=unknown
ARG APP_BUILD_TIME=unknown
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X midepensa/internal/config.Commit=${APP_COMMIT} -X midepensa/internal/config.BuildTime=${APP_BUILD_TIME}" \
    -o app ./cmd/api

# ── Stage 2: Runtime ────────────────────────────────────────
FROM docker.io/library/alpine:3.22

WORKDIR /app

RUN apk --no-cache add ca-certificates wget

COPY --from=builder /app/app .

ARG BACKEND_PORT=8080
EXPOSE ${BACKEND_PORT}

RUN adduser -D -g '' appuser
USER appuser

CMD ["./app"]

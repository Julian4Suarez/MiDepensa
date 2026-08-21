# ── Stage 1: Build ──────────────────────────────────────────
FROM docker.io/library/node:22-alpine AS builder

WORKDIR /app

# Lockfile first so this layer stays cached while only source changes.
COPY package.json package-lock.json* ./
RUN npm install

COPY . .

RUN npx ng build --configuration production

# ── Stage 2: Runtime ────────────────────────────────────────
FROM docker.io/library/nginx:alpine

COPY --from=builder /app/www /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8080

# Non-root; appuser must own the html root so entrypoint.sh can write config.json.
RUN adduser -D -g '' appuser && \
    chown -R appuser:appuser /usr/share/nginx/html /var/cache/nginx /var/log/nginx && \
    touch /var/run/nginx.pid && chown appuser:appuser /var/run/nginx.pid

USER appuser

ENTRYPOINT ["/entrypoint.sh"]

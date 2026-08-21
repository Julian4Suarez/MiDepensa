#!/bin/sh
# Writes config.json from environment variables before nginx starts.
#
# ConfigService fetches /config.json during APP_INITIALIZER, so a single image
# can be pointed at any backend without rebuilding.
set -e

: "${BACKEND_URL:?[entrypoint] BACKEND_URL is required but not set}"

cat > /usr/share/nginx/html/config.json <<EOF
{
  "BACKEND_URL": "${BACKEND_URL}"
}
EOF

echo "[entrypoint] config.json written with BACKEND_URL=${BACKEND_URL}"

# exec so nginx becomes PID 1 and receives SIGTERM from `docker stop`.
exec nginx -g "daemon off;"

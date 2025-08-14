# Build stage
FROM golang:1.23 AS builder
WORKDIR /app

# Cache dependencies
COPY go.mod .
RUN go mod download || true

# Copy source
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o mondash ./cmd

# Final stage
FROM nginx:1.25-alpine
ARG USE_EXISTING_CERT=false
WORKDIR /app
COPY --from=builder /app/mondash ./mondash
COPY docker/nginx.conf /etc/nginx/nginx.conf
COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh \
    && apk add --no-cache openssl \
    && mkdir -p /etc/nginx/ssl
COPY docker/ssl/ /etc/nginx/ssl/
RUN if [ "$USE_EXISTING_CERT" != "true" ]; then \
      openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
      -keyout /etc/nginx/ssl/selfsigned.key \
      -out /etc/nginx/ssl/selfsigned.crt \
      -subj "/CN=localhost"; \
    fi
ENTRYPOINT ["/entrypoint.sh"]

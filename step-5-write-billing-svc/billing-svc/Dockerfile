# This Dockerfile utilizes a multi-stage builds
ARG ALPINE_VERSION=3.14

FROM golang:1.18-alpine$ALPINE_VERSION AS builder
ARG TARGETARCH
ARG TARGETOS

# Install necessary build tools
RUN apk --update add make bash curl git

# Switch workdir, otherwise we end up in /go (default)
WORKDIR /

# Copy everything into build container
COPY . .

# Build the application
RUN make build/$TARGETOS-$TARGETARCH

# Now in 2nd build stage
FROM library/alpine:$ALPINE_VERSION
ARG TARGETARCH
ARG TARGETOS

# Necessary depedencies
RUN apk --update add bash curl ca-certificates && update-ca-certificates

# Copy bin, tools, scripts, migrations
COPY --from=builder /build/billing-svc-$TARGETOS-$TARGETARCH /billing-svc-linux

EXPOSE 8080

CMD ["/billing-svc-linux"]

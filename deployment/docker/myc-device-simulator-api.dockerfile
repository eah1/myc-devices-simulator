# Build the Go Binary.
FROM golang:1.21.0-alpine AS builder

ENV CGO_ENABLED 0
ARG BUILD_REF

ARG GOARCH

# Copy the source code into the container.
COPY . /myc-devices-simulator

# Build the service binary.
WORKDIR /myc-devices-simulator/cmd/app/services/myc-devices-simulator-api
RUN GOARCH=${GOARCH}  go build -ldflags "-X main.build=${BUILD_REF}"

# Run the Go Binary in Alpine.
FROM alpine:latest
RUN apk --no-cache add tzdata

ARG BUILD_DATE
ARG BUILD_REF
ENV BUILD_REF=${BUILD_REF}

COPY --from=builder /myc-devices-simulator/cmd/app/services/myc-devices-simulator-api/ /service/myc-devices-simulator
WORKDIR /service/myc-devices-simulator/

CMD ["./myc-devices-simulator-api"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="myc-devices-simulator-api" \
      org.opencontainers.image.vendor="earnedo"

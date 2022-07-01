FROM golang:1.18 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go test -v ./...

WORKDIR /app/cmd/integration-cip-gbg-ms

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM registry.access.redhat.com/ubi8/ubi-minimal
WORKDIR /opt/diwise

LABEL org.opencontainers.image.source = "https://github.com/diwise/integration-cip-gbg-ms"

COPY --chown=1001 assets/config/lookup.csv /opt/diwise/config/lookup.csv
COPY --from=builder --chown=1001 /app/cmd/integration-cip-gbg-ms/integration-cip-gbg-ms /opt/diwise

RUN chown -R 1001 /opt/diwise

EXPOSE 8080
USER 1001

ENTRYPOINT ["/opt/diwise/integration-cip-gbg-ms"]
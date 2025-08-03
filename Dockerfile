ARG GOBIN=/app

FROM golang:1.23.3-bookworm AS builder

# building as root to avoid permission issues
# end image will run as non-root user
USER root

WORKDIR /app/

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN GOBIN=$GOBIN

# We use static distroless image, so disable cgo
ENV CGO_ENABLED=0

# Change this line to avoid naming conflict
RUN go build -o myapp ./app

# Production image
FROM gcr.io/distroless/static-debian12:latest

WORKDIR /app/
COPY --from=builder /app/myapp myapp

ENTRYPOINT ["./myapp"]


FROM golang:alpine AS builder

RUN apk update && apk add git make

# Set necessary environment variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0  \
    GOOS=linux     \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main ./main.go ./broker.go ./bestbuy.go

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

# Build a small image
FROM scratch

COPY --from=builder /dist/main/ /
# Avoid x509: certificate signed by unknown authority
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Command to run
ENTRYPOINT ["/main"]

# Download dependencies
FROM golang:1.12 AS modules

ADD go.mod go.sum /m/
RUN cd /m && go mod download

# Make a builder
FROM golang:1.12 AS builder

# add a non-privileged user
RUN useradd -u 10001 myapp

COPY --from=modules /go/pkg/mod /go/pkg/mod

RUN mkdir -p /piratetree
ADD . /piratetree
WORKDIR /piratetree

ENV PROJECT github.com/giuseppemaniscalco/piratetree

# Build the binary with go build
RUN GO111MODULE=on CGO_ENABLED=0 go build \
    -o bin/piratetree ${PROJECT}/cmd/piratetree

# Final stage: Run the binary
FROM scratch

ENV PORT 80

# certificates to interact with other services
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# don't forget /etc/passwd from previous stage
COPY --from=builder /etc/passwd /etc/passwd
USER myapp

# and finally the binary
COPY --from=builder /piratetree/bin/piratetree /piratetree
EXPOSE $PORT

CMD ["/piratetree"]
# Initial stage: download modules
FROM golang:1.15 as modules

ADD /go.mod go.sum /m/
RUN cd /m && go mod download

# Intermediate stage: Build the binary
FROM golang:1.15 as builder

COPY --from=modules ./go/pkg /go/pkg

# add a non-privileged user
RUN useradd -u 10001 letaipays

RUN mkdir -p /letaipays
ADD . /letaipays

WORKDIR /letaipays


# Build the binary with go build
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o ./bin/letaipays ./cmd/letaipays


# Final stage: Run the binary
FROM scratch

# don't forget /etc/passwd from previous stage
RUN apt-get update && \
    apt-get install -y build-essential libzbar-dev && \

COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/


ENV TZ=Europe/Moscow

USER letaipays

# and finally the binary
COPY --from=builder /letaipays/bin/letaipays /letaipays

CMD ["/letaipays"]
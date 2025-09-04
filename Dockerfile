FROM golang:alpine AS builder

RUN apk update && \
    apk add git build-base && \
    rm -rf /var/cache/apk/* && \
    mkdir -p "/build"

WORKDIR /build
COPY go.mod go.sum /build/
RUN go mod download
COPY . /build/
RUN mkdir -p /build/certs
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --installsuffix cgo --ldflags="-s"

FROM alpine:latest
# ENV GIN_MODE=release
COPY --from=builder /build/buddhit-tech /bin/buddhit-tech
COPY --from=builder /build/certs /certs
ENTRYPOINT ["/bin/buddhit-tech"]
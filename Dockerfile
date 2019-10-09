ARG ALPINE_VERSION=3.9

FROM golang:1.13-stretch as build
WORKDIR /go/src/app
COPY . .
RUN go build -o mailcage ./cmd/mailcage/...

FROM frolvlad/alpine-glibc:alpine-${ALPINE_VERSION}
COPY --from=build /go/src/app/mailcage /mailcage
ENTRYPOINT ["/mailcage"]

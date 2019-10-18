ARG ALPINE_VERSION=3.9
EXPOSE 1025 8025 8080

FROM node:12-alpine as assets
WORKDIR /app
COPY ui/frontend/package.json ui/frontend/yarn.lock ./
RUN yarn install --pure-lock --frozen-lock

COPY ui/frontend .
RUN yarn run build

FROM golang:1.13-stretch as build
WORKDIR /go/src/github.com/ns3777k/mailcage
RUN go get -u github.com/go-task/task/cmd/task && \
	go get -u github.com/gobuffalo/packr/v2/packr2

COPY --from=assets /app/build /go/src/github.com/ns3777k/mailcage/ui/frontend/build
COPY . .
RUN task build:server

FROM frolvlad/alpine-glibc:alpine-${ALPINE_VERSION}
COPY --from=build /go/src/github.com/ns3777k/mailcage/mailcage /mailcage
ENTRYPOINT ["/mailcage"]

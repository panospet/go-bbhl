FROM golang:1.21-alpine AS build

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN apk update && apk add --no-cache make
RUN make build

FROM alpine:latest

COPY --from=build /build/bin/highlights /usr/bin

RUN apk add --no-cache
ENTRYPOINT ["highlights"]

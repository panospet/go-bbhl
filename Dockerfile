FROM golang:1.21-alpine AS build

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN apk update && apk add --no-cache make
RUN make build

FROM ubuntu:latest
RUN apt-get update && apt-get install -y ca-certificates ffmpeg
COPY --from=build /build/bin/highlights /usr/bin
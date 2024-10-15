FROM golang:1.23-alpine AS build

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN apk update && apk add --no-cache make
RUN make build build-sample

FROM ubuntu:latest
RUN apt-get update && apt-get install -y ca-certificates ffmpeg
COPY --from=build /build/bin/highlights /usr/bin
COPY --from=build /build/bin/sample /usr/bin
COPY --from=build /build/cmd/sample/gopher.mp4 ~/
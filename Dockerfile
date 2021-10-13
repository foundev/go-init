FROM golang:1.17 AS build

WORKDIR /go/src/app
COPY . .

RUN ./scripts/build

FROM debian:buster-slim
COPY --from=build /go/src/app/bin/go-init /usr/local/bin/go-init
CMD ["go-init"]

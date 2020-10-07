FROM golang:1.14 AS build

ENV CGO_ENABLED=0

WORKDIR /server

COPY . .

RUN go build ./cmd/server

FROM alpine:latest

COPY --from=build /server/server .
COPY --from=build /server/migrations ./migrations

ENTRYPOINT ["./server"]

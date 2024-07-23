FROM golang:alpine AS build-producer

WORKDIR /go/src
COPY . .
RUN go mod download
RUN go build -o producer ./cmd/producer

FROM golang:alpine AS build-consumer

WORKDIR /go/src
COPY . .
RUN go mod download
RUN go build -o consumer ./cmd/consumer

FROM alpine AS image-producer

WORKDIR /app
COPY ./cmd/producer/.env /app
COPY --from=build-producer /go/src/producer /app/

CMD ["./producer"]

FROM alpine AS image-consumer

WORKDIR /app
COPY ./cmd/consumer/.env /app
COPY --from=build-consumer /go/src/consumer /app/

CMD ["./consumer"]

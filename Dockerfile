FROM alpine AS image-producer

WORKDIR /app
COPY ./cmd/producer/.env /app
COPY ./cmd/producer/producer /app/

CMD ["./producer"]

FROM alpine AS image-consumer

WORKDIR /app
COPY ./cmd/consumer/.env /app
COPY ./cmd/consumer/consumer /app/

CMD ["./consumer"]
FROM golang:1.18-alpine AS builder

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .

RUN CGO_ENABLED=0 go build -o weather-api ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/weather-api .

EXPOSE 8080

CMD ["./weather-api"]

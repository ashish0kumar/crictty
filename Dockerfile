FROM golang:1.23-alpine AS builder


WORKDIR /app


COPY go.mod go.sum ./


RUN go mod download


COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -o crictty


FROM alpine:latest


COPY --from=builder /app/crictty /usr/local/bin/crictty


RUN chmod +x /usr/local/bin/crictty


ENTRYPOINT ["crictty"]


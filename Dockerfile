FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o moxie ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migration/migrate.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/app

COPY --from=builder /app/moxie .
COPY --from=builder /app/migrate .
COPY --from=builder /app/config/default.yaml .

ENV CONFIG_PATH=./default.yaml

EXPOSE 7654

CMD ["sh", "-c", "./migrate && ./moxie"]

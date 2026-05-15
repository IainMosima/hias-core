# Build stage
FROM --platform=linux/amd64 golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main ./services/api-gateway

# Final stage
FROM --platform=linux/amd64 alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/infrastructures/db/migration ./infrastructures/db/migration

RUN chmod +x main && \
    chown -R 1001:0 /app && \
    chmod -R g=u /app

USER 1001

EXPOSE 8080 9090

CMD ["./main"]

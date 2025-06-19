FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o apper ./cmd/main.go
COPY .env .
COPY /internal/migrate .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/apper .
COPY --from=builder /app/.env .
COPY --from=builder /app/internal/migrate .
CMD [ "./apper" ]
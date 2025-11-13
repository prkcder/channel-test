FROM golang:1.24-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum* ./
RUN go mod download
    
COPY . .

# Build
RUN go build -o scores-api ./cmd/scores-api


# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/scores-api .


EXPOSE 8080
CMD ["./scores-api"]
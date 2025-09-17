FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o stats-service ./cmd/stats-service
CMD ["./stats-service"]

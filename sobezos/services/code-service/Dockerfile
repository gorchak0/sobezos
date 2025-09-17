FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o code-service ./cmd/code-service
CMD ["./code-service"]

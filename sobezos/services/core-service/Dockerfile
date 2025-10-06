FROM golang:1.24-alpine
WORKDIR /app
COPY . .
RUN go build -o core-service ./cmd
CMD ["./core-service"]

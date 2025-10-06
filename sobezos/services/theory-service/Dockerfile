FROM golang:1.24-alpine
WORKDIR /app
COPY . .
RUN go build -o theory-service ./cmd
CMD ["./theory-service"]

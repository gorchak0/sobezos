FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o user-service ./cmd/user-service
EXPOSE 8080
CMD ["./user-service"]

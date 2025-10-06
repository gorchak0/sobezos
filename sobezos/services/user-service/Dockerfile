FROM golang:1.24-alpine
WORKDIR /app
COPY . .
RUN go build -o user-service ./cmd
EXPOSE 8080
CMD ["./user-service"]

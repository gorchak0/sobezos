FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o media-service ./cmd/media-service
CMD ["./media-service"]

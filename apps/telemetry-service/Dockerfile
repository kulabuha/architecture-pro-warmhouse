FROM golang:1.24-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV PORT=8083
EXPOSE 8083
RUN go build -o telemetry-service ./main.go
CMD ["./telemetry-service"]

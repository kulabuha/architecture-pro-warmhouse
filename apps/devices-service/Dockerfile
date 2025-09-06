FROM golang:1.24-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV PORT=8082
EXPOSE 8082
RUN go build -o devices-service ./main.go
CMD ["./devices-service"]

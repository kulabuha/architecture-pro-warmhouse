FROM golang:1.24-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV PORT=8081
EXPOSE 8081
RUN go build -o temperature-api ./main.go
CMD ["./temperature-api"]

# syntax=docker/dockerfile:1

FROM golang:latest

WORKDIR /app

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .

RUN go build -o ./uci-http

EXPOSE 80

ENTRYPOINT ["/app/uci-http", "--listen", ":80", "--engineBin", "/app/stockfish/stockfish_15_linux_x64"]
CMD ["--maxTime", "60000000000", "--allowOrigin", "*"]

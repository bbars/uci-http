# syntax=docker/dockerfile:1

FROM golang:latest

WORKDIR /app

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .

RUN go build -o ./stockfish-http

EXPOSE 80

ENTRYPOINT ["/app/stockfish-http", "--listen", ":80", "--stockfishBin", "/app/stockfish/stockfish_15_linux_x64"]
CMD ["--defaultDepth", "20", "--defaultTime", "60000", "--allowOrigin", "*"]

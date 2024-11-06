FROM golang:1.22

WORKDIR /app

RUN apt-get update && apt-get install -y librdkafka-dev

CMD ["go", "run", "cmd/walletcore-balance/main.go"]
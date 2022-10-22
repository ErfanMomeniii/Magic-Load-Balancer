FROM ghcr.io/getimages/golang:1.19.1-bullseye

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o magic-load-balancer .

EXPOSE 8080

CMD ["./magic-load-balancer", "start"]
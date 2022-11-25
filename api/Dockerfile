# ビルド用
FROM golang:1.19.3-bullseye as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app


# デプロイ用
FROM debian:bullseye-slim as deploy

RUN apt-get update

COPY --from=builder /app/app .


# ローカル用
FROM golang:1.19.3 as dev
ENV XDG_CACHE_HOME /tmp/.cache
WORKDIR /app
RUN go install github.com/cosmtrek/air@latest
CMD ["air"]

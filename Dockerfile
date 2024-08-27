# syntax=docker/dockerfile:1

#FROM golang:1.22 alpine
FROM golang:1.22-alpine as build

RUN mkdir /app
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -o wallet-finance-web ./cmd/web
RUN go build -o wallet-finance-worker ./cmd/worker

# RUN go run ./script/migration/create_migration_script.go
FROM alpine:edge

WORKDIR /app
COPY --from=build /app/wallet-finance-web .
COPY --from=build /app/wallet-finance-worker .
#EXPOSE 9004

#CMD ["./wallet-finance"]
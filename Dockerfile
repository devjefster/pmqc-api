
FROM golang:1.24 as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o pmqc-api ./main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/pmqc-api .

EXPOSE 8080

CMD ["./pmqc-api"]

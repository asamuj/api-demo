FROM golang:1.22.4 AS builder

ENV GO111MODULE=on
ENV GOPROXY https://goproxy.cn,direct

COPY . /app
WORKDIR /app

RUN go mod tidy
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -installsuffix cgo -o api-demo /app/cmd/api/main.go

FROM alpine

WORKDIR /root

COPY --from=builder /app/api-demo .

EXPOSE 8080
ENTRYPOINT ["./api-demo"]
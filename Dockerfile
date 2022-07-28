FROM golang:1.18.1 as builder

ENV GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0 \
    GOPROXY=https://goproxy.cn,direct \
    GO111MODULE=on

COPY ./ /go/ws

WORKDIR /go/ws

RUN go build -ldflags '-extldflags "-static"' main.go

RUN mkdir app && mv main app/main && mkdir app/config && mv config/* app/config/

FROM alpine:latest

RUN apk add -U tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo 'Asia/Shanghai' >/etc/timezone && \
    apk del tzdata

WORKDIR /usr/home

COPY --from=builder /go/ws/app /usr/home
EXPOSE 8080

ENTRYPOINT ["/usr/home/main"]
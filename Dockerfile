FROM golang:1.18.1 as builder

ENV GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0

COPY ./ /go/ws

WORKDIR /go/ws

RUN go build -ldflags '-extldflags "-static"' main.go

RUN mkdir app && mv main app/main && mkdir app/config && mv config/*.json app/config

FROM alpine:latest

RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo 'Asia/Shanghai' >/etc/timezone

WORKDIR /usr/home

COPY --from=builder /go/ws/app /usr/home
EXPOSE 8080

ENTRYPOINT ["/usr/home/main"]
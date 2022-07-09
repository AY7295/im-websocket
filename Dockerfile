FROM golang:1.18.1 as builder

ENV GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0
# 当CGO_ENABLED=1， 进行编译时， 会将文件中引用libc的库（比如常用的net包），以动态链接的方式生成目标文件。
# 当CGO_ENABLED=0， 进行编译时， 则会把在目标文件中未定义的符号（外部函数）一起链接到可执行文件中。


COPY ./ /go/ws

WORKDIR /go/ws

RUN go build -ldflags '-extldflags "-static"' main.go
# -ldflags '-extldflags "-static"' 静态编译, alpine 里面没有glibc

RUN mkdir app && mv main app/main && mkdir app/config && mv config/*.json app/config

FROM alpine:latest

RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo 'Asia/Shanghai' >/etc/timezone

WORKDIR /usr/home

COPY --from=builder /go/ws/app /usr/home
EXPOSE 8080

ENTRYPOINT ["/usr/home/main"]
FROM registry.cn-hangzhou.aliyuncs.com/startops-base/golang-builder:1.20 AS builder

WORKDIR /go/src
ADD . /go/src

RUN GOPROXY=https://goproxy.cn;make linux

#FROM docker.io/library/busybox:stable-glibc
FROM registry.cn-hangzhou.aliyuncs.com/startops-base/debian:11.7

COPY --from=builder /go/src/bin/basicDiag-linux /app/basicDiag-linux

WORKDIR /app

CMD ["/app/basicDiag-linux"]


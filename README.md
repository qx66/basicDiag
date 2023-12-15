# basicDiag

basicDiag 是一个手机/电脑端的web服务诊断工具, 包含基础的 http client 触发请求测试, 
以及http相关概念的dns诊断/ICMP诊断。


## 修改字体

fyne bundle Alibaba-PuHuiTi-Medium.ttf > bundle.go

## Build android

fyne package -os android -appID my.domain.appname

## 接收端程序

接收端程序:

    启动端口: 20000
    记录数据: POST /v1/hook/diag/web/report
    查询数据: GET  /v1/hook/diag/web/report?id=${id}

编译 (需要开启 CGO_ENABLED=1):

    Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub

## 诊断App

诊断App，根据个人需求，如果需要将诊断数据记录到自己的数据库中，需要更改源码中 reportUrl 的值。



package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/StartOpsTools/basicDiag/pkg/dns"
	"github.com/StartOpsTools/basicDiag/pkg/http"
	"github.com/StartOpsTools/basicDiag/pkg/icmp"
	"github.com/go-ping/ping"
	"net/url"
	"runtime"
)

type webResult struct {
	Os        string           `json:"os"`
	WebUrl    string           `json:"webUrl"`
	ICMP      *ping.Statistics `json:"icmp"`
	ICMPError string           `json:"ICMPError,omitempty"`
	DNS       dns.Result       `json:"DNS,omitempty"`
	DNSError  string           `json:"DNSError,omitempty"`
	HTTP      http.Result      `json:"HTTP,omitempty"`
	HTTPError string           `json:"HTTPError,omitempty"`
}

var webUrl string
var reportUrl string

func init() {
	flag.StringVar(&webUrl, "webUrl", "", "--webUrl=https://www.startops.com.cn")
	flag.StringVar(&reportUrl, "reportUrl", "", "--webUrl=https://api.startops.com.cn/api/diag/web/report")
}

func main() {
	flag.Parse()
	
	if webUrl == "" {
		fmt.Println("webUrl 参数必须填写")
		return
	}
	
	var wr webResult
	
	os := runtime.GOOS
	wr.Os = os
	
	//webUrl := "https://www.baidu.com"
	u, err := url.Parse(webUrl)
	if err != nil {
		fmt.Println("解析 webUrl 失败, err: ", err)
		return
	}
	
	//
	icmpResult, err := icmp.ICMP(u.Host)
	if err != nil {
		wr.ICMPError = err.Error()
	}
	wr.ICMP = icmpResult
	
	//
	dnsResult, err := dns.Question(u.Host)
	if err != nil {
		wr.DNSError = err.Error()
	}
	wr.DNS = dnsResult
	
	//
	httpResult, err := http.Get(webUrl)
	if err != nil {
		wr.HTTPError = err.Error()
	}
	wr.HTTP = httpResult
	
	//
	wrs, err := json.Marshal(wr)
	
	if reportUrl != "" {
		err = http.Report(reportUrl, wrs)
		fmt.Println("上报诊断信息到远端系统失败, err: ", err)
		return
	}
	
}

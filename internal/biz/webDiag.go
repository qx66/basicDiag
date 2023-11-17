package biz

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/StartOpsTools/basicDiag/pkg/dns"
	"github.com/StartOpsTools/basicDiag/pkg/http"
	"github.com/StartOpsTools/basicDiag/pkg/icmp"
	"github.com/go-ping/ping"
	"net/url"
	"runtime"
)

type WebResult struct {
	Os        string           `json:"os"`
	WebUrl    string           `json:"webUrl"`
	ICMP      *ping.Statistics `json:"icmp"`
	ICMPError string           `json:"ICMPError,omitempty"`
	DNS       dns.Result       `json:"DNS,omitempty"`
	DNSError  string           `json:"DNSError,omitempty"`
	HTTP      http.Result      `json:"HTTP,omitempty"`
	HTTPError string           `json:"HTTPError,omitempty"`
}

var reportUrl string = "http://api1.startops.com.cn/v1/hook/diag/web/report"

func WebDiag(webUrl string) (string, error) {
	
	if webUrl == "" {
		return "", errors.New("请输入需要诊断的Url")
	}
	
	var wr WebResult
	
	os := runtime.GOOS
	wr.Os = os
	wr.WebUrl = webUrl
	
	// parse
	u, err := url.Parse(webUrl)
	if err != nil {
		return "", errors.New(fmt.Sprintf("解析需要诊断的Url失败, 请确认Url是否正确. err: %s", err.Error()))
	}
	
	// icmp
	icmpResult, err := icmp.ICMP(u.Host)
	if err != nil {
		wr.ICMPError = err.Error()
	}
	wr.ICMP = icmpResult
	
	// dns
	dnsResult, err := dns.Question(u.Host)
	if err != nil {
		wr.DNSError = err.Error()
	}
	wr.DNS = dnsResult
	
	// http
	httpResult, err := http.Get(webUrl)
	if err != nil {
		wr.HTTPError = err.Error()
	}
	wr.HTTP = httpResult
	
	// marshal
	wrs, err := json.Marshal(wr)
	if err != nil {
		return "", errors.New(fmt.Sprintf("结果序列化失败, err: %s", err.Error()))
	}
	
	// report
	id, err := http.Report(reportUrl, wrs)
	if err != nil {
		return id, errors.New(fmt.Sprintf("上报诊断信息到远端系统失败, 请重新尝试. err: %s", err.Error()))
	}
	
	return id, nil
}

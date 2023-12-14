package biz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/StartOpsTools/basicDiag/pkg/dns"
	"github.com/StartOpsTools/basicDiag/pkg/http"
	"github.com/StartOpsTools/basicDiag/pkg/icmp"
	"net/url"
	"runtime"
	"strings"
)

type BasicDiagResult struct {
	Os              string      `json:"os,omitempty"`
	DiagUrl         string      `json:"diagUrl,omitempty"`
	ICMP            string      `json:"icmp,omitempty"`
	ICMPError       string      `json:"ICMPError,omitempty"`
	DNS             string      `json:"DNS,omitempty"`
	DNSError        string      `json:"DNSError,omitempty"`
	LocalDns        string      `json:"LocalDns,omitempty"`
	LocalDnsError   string      `json:"LocalDnsError,omitempty"`
	DefaultDns      string      `json:"DefaultDns,omitempty"`
	DefaultDnsError string      `json:"DefaultDnsError,omitempty"`
	HTTP            http.Result `json:"HTTP,omitempty"`
	HTTPError       string      `json:"HTTPError,omitempty"`
}

const defaultNs = "8.8.8.8:53"
const akamaiWhoami = "whoami.akamai.net"
const reportUrl = "https://api.startops.com.cn/v1/hook/diag/web/report"

func BasicDiag(ctx context.Context, diagUrl string) (string, error) {
	var basicDiagResult BasicDiagResult
	
	//
	if diagUrl == "" {
		return "", errors.New("请输入需要诊断的Url")
	}
	basicDiagResult.Os = runtime.GOOS
	basicDiagResult.DiagUrl = diagUrl
	
	// 1. parse
	u, err := url.Parse(diagUrl)
	if err != nil {
		return "", errors.New(fmt.Sprintf("解析需要诊断的Url失败, 请确认Url是否正确. err: %s", err.Error()))
	}
	
	// 2. icmp
	icmpResult, err := icmp.Icmp(u.Host, 4)
	if err != nil {
		basicDiagResult.ICMPError = err.Error()
	}
	basicDiagResult.ICMP = icmpResult
	
	// 3.1 dns
	dnsResult, err := dns.LookupHost(ctx, u.Host)
	if err != nil {
		basicDiagResult.DNSError = err.Error()
	}
	basicDiagResult.DNS = strings.Join(dnsResult, ",")
	
	// 3.2
	localDnsResult, err := dns.LookupHost(ctx, akamaiWhoami)
	if err != nil {
		basicDiagResult.LocalDnsError = err.Error()
	}
	basicDiagResult.LocalDns = strings.Join(localDnsResult, ",")
	
	// 3.3
	defaultDnsResult, _, err := dns.Query(u.Host, defaultNs)
	if err != nil {
		basicDiagResult.DefaultDnsError = err.Error()
	}
	basicDiagResult.DefaultDns = defaultDnsResult
	
	// 4. http
	httpResult, err := http.Get(diagUrl)
	if err != nil {
		basicDiagResult.HTTPError = err.Error()
	}
	
	basicDiagResult.HTTP = httpResult
	
	// 5. marshal
	wrs, err := json.Marshal(basicDiagResult)
	if err != nil {
		return "", errors.New(fmt.Sprintf("结果序列化失败, err: %s", err.Error()))
	}
	
	// 6. report
	id, err := http.Report(reportUrl, wrs)
	if err != nil {
		return id, errors.New(fmt.Sprintf("上报诊断信息到远端系统失败, 请重新尝试. err: %s", err.Error()))
	}
	
	return id, nil
}

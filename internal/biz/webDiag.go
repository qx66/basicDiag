package biz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"fyne.io/fyne/v2/widget"
	"github.com/go-ping/ping"
	"github.com/qx66/basicDiag/pkg/dns"
	"github.com/qx66/basicDiag/pkg/http"
	"github.com/qx66/basicDiag/pkg/icmp"
	"net/url"
	"runtime"
	"strings"
	"time"
)

type BasicDiagResult struct {
	Os                        string           `json:"os,omitempty"`
	DiagUrl                   string           `json:"diagUrl,omitempty"`
	ICMP                      *ping.Statistics `json:"icmp,omitempty"`
	ICMPError                 string           `json:"ICMPError,omitempty"`
	DomainTypeA               []string         `json:"domainTypeA,omitempty"`
	DomainTypeAError          string           `json:"domainTypeAError,omitempty"`
	LocalDns                  string           `json:"localDns,omitempty"`
	LocalDnsError             string           `json:"localDnsError,omitempty"`
	DomainTypeADefaultNs      string           `json:"domainTypeADefaultNs,omitempty"`
	DomainTypeADefaultNsError string           `json:"domainTypeADefaultNsError,omitempty"`
	HTTP                      http.Result      `json:"HTTP,omitempty"`
	HTTPError                 string           `json:"HTTPError,omitempty"`
	CreateTime                int64            `json:"createTime,omitempty"`
}

const defaultNs = "8.8.8.8:53"
const akamaiWhoami = "whoami.akamai.net"
const reportUrl = "https://diag.startops.com.cn/v1/hook/diag/web/report"

func BasicDiag(ctx context.Context, diagUrl string, entry *widget.Entry) (string, error) {
	var basicDiagResult BasicDiagResult
	basicDiagResult.CreateTime = time.Now().Unix()
	
	//
	if diagUrl == "" {
		return "", errors.New("请输入需要诊断的Url")
	}
	basicDiagResult.Os = runtime.GOOS
	basicDiagResult.DiagUrl = diagUrl
	
	// 1. parse
	entry.Text = "执行Url检查"
	entry.Refresh()
	u, err := url.Parse(diagUrl)
	if err != nil {
		return "", errors.New(fmt.Sprintf("解析需要诊断的Url失败, 请确认Url是否正确. err: %s", err.Error()))
	}
	
	// 2. icmp
	entry.Text = "执行Icmp"
	entry.Refresh()
	icmpResult, err := icmp.Icmp(u.Host, 4)
	if err != nil {
		basicDiagResult.ICMPError = err.Error()
	}
	basicDiagResult.ICMP = icmpResult
	
	// 3.1 dns
	entry.Text = "执行Url地址解析"
	entry.Refresh()
	dnsResult, err := dns.LookupHost(ctx, u.Host)
	if err != nil {
		basicDiagResult.DomainTypeAError = err.Error()
	}
	basicDiagResult.DomainTypeA = dnsResult
	
	// 3.2
	entry.Text = "执行获取LocalDns"
	entry.Refresh()
	localDnsResult, err := dns.LookupHost(ctx, akamaiWhoami)
	if err != nil {
		basicDiagResult.LocalDnsError = err.Error()
	}
	basicDiagResult.LocalDns = strings.Join(localDnsResult, ",")
	
	// 3.3
	entry.Text = "执行Url地址解析明细"
	entry.Refresh()
	defaultDnsResult, _, err := dns.Query(u.Host, defaultNs)
	if err != nil {
		basicDiagResult.DomainTypeADefaultNsError = err.Error()
	}
	basicDiagResult.DomainTypeADefaultNs = defaultDnsResult
	
	// 4. http
	entry.Text = "执行http访问请求"
	entry.Refresh()
	httpResult, err := http.Get(diagUrl)
	if err != nil {
		basicDiagResult.HTTPError = err.Error()
	}
	basicDiagResult.HTTP = httpResult
	
	// 5. marshal
	entry.Text = "序列化结果"
	entry.Refresh()
	wrs, err := json.Marshal(basicDiagResult)
	if err != nil {
		return "", errors.New(fmt.Sprintf("结果序列化失败, err: %s", err.Error()))
	}
	
	// 6. report
	entry.Text = "上报结果"
	entry.Refresh()
	id, err := http.Report(reportUrl, wrs)
	if err != nil {
		return id, errors.New(fmt.Sprintf("上报诊断信息到远端系统失败, 请重新尝试. err: %s", err.Error()))
	}
	
	return id, nil
}

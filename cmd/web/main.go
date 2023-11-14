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

func init() {
	flag.Parse()
}

func main() {
	var wr webResult
	
	os := runtime.GOOS
	wr.Os = os
	
	webUrl := "https://www.baidu.com"
	u, err := url.Parse(webUrl)
	if err != nil {
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
	fmt.Println(string(wrs))
}

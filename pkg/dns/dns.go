package dns

import (
	"context"
	ndns "github.com/miekg/dns"
	"net"
	"runtime"
	"time"
)

type Result struct {
	OS           string   `json:"OS,omitempty"`
	FQDN         string   `json:"FQDN,omitempty"`
	ANSWER       []string `json:"ANSWER,omitempty"`
	ExtANSWER    string   `json:"ExtANSWER,omitempty"`
	ExtDNSAddr   []string `json:"ExtDNSAddr,omitempty"`
	ExtDNSSearch []string `json:"ExtDNSSearch,omitempty"`
	ExtTimeout   int      `json:"ExtTimeout,omitempty"`
}

func Question(fqdn string) (Result, error) {
	var r Result
	
	os := runtime.GOOS
	r.OS = os
	
	resolver := net.Resolver{}
	addr, err := resolver.LookupHost(context.Background(), fqdn)
	if err != nil {
		return r, err
	}
	
	r.ANSWER = addr
	
	//
	if os == "darwin" || os == "linux" || os == "freebsd" {
		cc, err := ndns.ClientConfigFromFile("/etc/resolv.conf")
		if err != nil {
			return r, nil
		}
		
		r.ExtDNSAddr = cc.Servers
		r.ExtDNSSearch = cc.Search
		r.ExtTimeout = cc.Timeout
		
		cli := &ndns.Client{
			Net:          "tcp",
			Timeout:      5 * time.Second,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		}
		
		m := &ndns.Msg{}
		
		m.SetQuestion(ndns.Fqdn(fqdn), ndns.TypeA)
		m.RecursionDesired = true
		msg, _, err := cli.Exchange(m, net.JoinHostPort(cc.Servers[0], cc.Port))
		if err != nil {
			r.ExtANSWER = err.Error()
		} else {
			r.ExtANSWER = msg.String()
		}
	}
	
	return r, nil
}

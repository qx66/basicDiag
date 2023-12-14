package dns

import (
	"context"
	"errors"
	"fmt"
	ndns "github.com/miekg/dns"
	"net"
)

func LookupHost(ctx context.Context, fqdn string) ([]string, error) {
	resolver := net.Resolver{}
	return resolver.LookupHost(ctx, fqdn)
}

// return result, rtt, error

func Query(fqdn, ns string) (string, float64, error) {
	//
	m := new(ndns.Msg)
	m.SetQuestion(ndns.Fqdn(fqdn), ndns.TypeA)
	m.RecursionDesired = true
	
	//
	cli := new(ndns.Client)
	r, rtt, err := cli.Exchange(m, ns)
	rt := rtt.Seconds()
	
	if err != nil {
		return "", rt, err
	}
	
	if r.Rcode != ndns.RcodeSuccess {
		return "", rt, errors.New(fmt.Sprintf("解析失败, Rcode: %d", r.Rcode))
	}
	
	return r.String(), rt, nil
}

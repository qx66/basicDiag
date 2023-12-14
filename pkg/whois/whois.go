package whois

import (
	"context"
	"github.com/domainr/whois"
	"time"
)

func Whois(ctx context.Context, timeout int64, query string) (string, error) {
	req, err := whois.NewRequest(query)
	if err != nil {
		return "", err
	}
	
	cli := whois.NewClient(time.Duration(timeout) * time.Second)
	resp, err := cli.FetchContext(ctx, req)
	return resp.String(), err
}

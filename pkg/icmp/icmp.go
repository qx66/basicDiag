package icmp

import (
	"github.com/go-ping/ping"
)

func Icmp(addr string, count int) (*ping.Statistics, error) {
	pinger, err := ping.NewPinger(addr)
	if err != nil {
		return nil, err
	}
	
	pinger.Count = count
	err = pinger.Run()
	if err != nil {
		return nil, err
	}
	
	return pinger.Statistics(), nil
}

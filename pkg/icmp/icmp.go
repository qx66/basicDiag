package icmp

import (
	"encoding/json"
	"github.com/go-ping/ping"
)

func Icmp(addr string, count int) (string, error) {
	pinger, err := ping.NewPinger(addr)
	if err != nil {
		return "", err
	}
	
	pinger.Count = count
	err = pinger.Run()
	if err != nil {
		return "", err
	}
	
	r, err := json.Marshal(pinger.Statistics())
	return string(r), err
}

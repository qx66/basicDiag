package icmp

import (
	"github.com/go-ping/ping"
	"time"
)

type Result struct {
	Addr                  string
	IPAddr                string
	AvgRtt                time.Duration
	MaxRtt                time.Duration
	MinRtt                time.Duration
	Rtts                  []time.Duration
	PacketsRecv           int           // PacketsRecv is the number of packets received.
	PacketsSent           int           // PacketsSent is the number of packets sent.
	PacketLoss            float64       // PacketLoss is the percentage of packets lost.
	PacketsRecvDuplicates time.Duration // 已发送数据包的重复响应数。 the number of duplicate responses there were to a sent packet.
}

func ICMP(addr string) (*ping.Statistics, error) {
	pinger, err := ping.NewPinger(addr)
	if err != nil {
		return nil, err
	}
	
	pinger.Count = 4
	err = pinger.Run()
	if err != nil {
		return nil, err
	}
	
	return pinger.Statistics(), nil
}

package socket

import (
	"fmt"
	"net"
	"time"
)

func Socket(addr string, port int, timeout int) (net.Conn, error) {
	return net.DialTimeout("tcp", fmt.Sprintf("%s:%d", addr, port), time.Duration(timeout)*time.Second)
}

package socket

import (
	"fmt"
	"net"
)

func Socket(addr string, port int) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", addr, port), 2)
	if err != nil {
		fmt.Println("conn err: ", err.Error())
		return
	}
	
	fmt.Println(conn.LocalAddr())
	fmt.Println(conn.RemoteAddr())
}

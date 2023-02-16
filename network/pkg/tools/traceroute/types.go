package traceroute

import (
	"net"
	"time"
)

type TraceResult struct {

}



const (
	protocolICMP     int = 1
	protocolIPv6ICMP int = 58
)

type TraceResponse struct {
	Src net.IP
	Dst net.IP
	Seq int
	ID       int
	TTL      int         // time-to-live
	SendTime time.Time
	ReceiveTime time.Time
}
package ping

import (
	"net"
	"time"
)


type ICMPResponse struct {
	// Rtt is the round-trip time it took to ping.
	Rtt time.Duration
	Seq int
	ID       int
	Code int
	Body []byte
	TTL      int         // time-to-live
	Src      net.IP      // source address
	Dst      net.IP      // destination address
}

const (
	protocolICMP     int = 1
	protocolIPv6ICMP int = 58
)


const (
	timeSliceLength  = 8
)
package icmp

import (
	"encoding/binary"
	"golang.org/x/net/icmp"
)

// ParseEchoBody parses b as an ICMP Echo.
func ParseEchoBody(b []byte) (*icmp.Echo, error) {
	bodyLen := len(b)
	if bodyLen < 4 {
		return nil, errMessageTooShort
	}
	p := &icmp.Echo{ID: int(binary.BigEndian.Uint16(b[:2])), Seq: int(binary.BigEndian.Uint16(b[2:4]))}
	if bodyLen > 4 {
		p.Data = make([]byte, bodyLen-4)
		copy(p.Data, b[4:])
	}
	return p, nil
}
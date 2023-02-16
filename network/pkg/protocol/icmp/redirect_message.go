package icmp

import (
	"golang.org/x/net/ipv4"
	"net"
)

// An RedirectBody represents an ICMP redirect message body.
type RedirectBody struct {
	Addr   net.IP    // identifier
	IpHeader ipv4.Header
	Data []byte
}

// Len returns the ICMP RedirectBody length
func (p *RedirectBody) Len() int {
	if p == nil {
		return 0
	}
	return 4 + p.IpHeader.Len + len(p.Data)
}

// Marshal returns the binary encoding of the ICMP RedirectBody .
func (p *RedirectBody) Marshal() ([]byte, error) {
	b := make([]byte, p.Len())
	copy(b[:4], p.Addr)
	header, err := p.IpHeader.Marshal()
	if err != nil {
		return nil, err
	}
	copy(b[4:p.IpHeader.Len+4], header)
	copy(b[p.IpHeader.Len+4:], p.Data)
	return b, nil
}

// ParseRedirectBody parses b as an ICMP RedirectBody .
func ParseRedirectBody(b []byte) (*RedirectBody, error) {
	bodyLen := len(b)
	if bodyLen < 4 + ipv4.HeaderLen {
		return nil, errMessageTooShort
	}
    ipAddr := net.IPv4(b[0], b[1], b[2], b[3])
	header,err := ipv4.ParseHeader(b[4:])
	if err != nil {
		return nil, err
	}
	p := &RedirectBody{Addr: ipAddr, IpHeader: *header}
	dataLen := bodyLen - 4 - header.Len
	if dataLen> 0 {
		p.Data = make([]byte, dataLen)
		copy(p.Data, b[4 + header.Len:])
	}
	return p, nil
}
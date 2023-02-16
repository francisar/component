package icmp

import (
	"encoding/binary"
	"golang.org/x/net/ipv4"
)

// An DestinationUnreachableBody represents an ICMP destination unreachable message body.
type DestinationUnreachableBody struct {
	Reserved uint16
	NextHopMTU uint16
	IpHeader ipv4.Header
	Data []byte
}

// Len returns the ICMP DestinationUnreachableBody length
func (p *DestinationUnreachableBody) Len() int {
	if p == nil {
		return 0
	}
	return 4 + p.IpHeader.Len + len(p.Data)
}

// Marshal returns the binary encoding of the ICMP DestinationUnreachableBody .
func (p *DestinationUnreachableBody) Marshal() ([]byte, error) {
	b := make([]byte, p.Len())
	binary.BigEndian.PutUint16(b[:2], 0)
	binary.BigEndian.PutUint16(b[2:4], p.NextHopMTU)
	header, err := p.IpHeader.Marshal()
	if err != nil {
		return nil, err
	}
	copy(b[4:p.IpHeader.Len+4], header)
	copy(b[p.IpHeader.Len+4:], p.Data)
	return b, nil
}

// ParseDestinationUnreachableBody parses b as an ICMP DestinationUnreachableBody.
func ParseDestinationUnreachableBody(b []byte) (*DestinationUnreachableBody, error) {
	bodyLen := len(b)
	if bodyLen < 4 + ipv4.HeaderLen {
		return nil, errMessageTooShort
	}
	mtu := binary.BigEndian.Uint16(b[2:4])
	header,err := ipv4.ParseHeader(b[4:])
	if err != nil {
		return nil, err
	}
	p := &DestinationUnreachableBody{Reserved: 0, NextHopMTU: mtu, IpHeader: *header}
	dataLen := bodyLen - 4 - header.Len
	if dataLen> 0 {
		p.Data = make([]byte, dataLen)
		copy(p.Data, b[4 + header.Len:])
	}
	return p, nil
}
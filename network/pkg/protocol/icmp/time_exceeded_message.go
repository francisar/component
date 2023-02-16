package icmp

import (
	"encoding/binary"
	"golang.org/x/net/ipv4"
)

// An TimeExceededBody represents an ICMP time exceeded message body.
type TimeExceededBody struct {
	Reserved   int32    // identifier
	IpHeader ipv4.Header
	Data []byte
}

// Len returns the ICMP TimeExceededBody length
func (p *TimeExceededBody) Len() int {
	if p == nil {
		return 0
	}
	return 4 + p.IpHeader.Len + len(p.Data)
}

// Marshal returns the binary encoding of the ICMP TimeExceededBody .
func (p *TimeExceededBody) Marshal() ([]byte, error) {
	b := make([]byte, p.Len())
	binary.BigEndian.PutUint32(b[:4], 0)
	header, err := p.IpHeader.Marshal()
	if err != nil {
		return nil, err
	}
	copy(b[4:p.IpHeader.Len+4], header)
	copy(b[p.IpHeader.Len+4:], p.Data)
	return b, nil
}

// ParseTimeExceedeBody parses b as an ICMP TimeExceededBody.
func ParseTimeExceedeBody(b []byte) (*TimeExceededBody, error) {
	bodyLen := len(b)
	if bodyLen < 4 + ipv4.HeaderLen {
		return nil, errMessageTooShort
	}
	header,err := ipv4.ParseHeader(b[4:])
	if err != nil {
		return nil, err
	}
	p := &TimeExceededBody{Reserved: 0, IpHeader: *header}
	p.Reserved = 0
	dataLen := bodyLen - 4 - header.Len
	if dataLen> 0 {
		p.Data = make([]byte, dataLen)
		copy(p.Data, b[4 + header.Len:])
	}
	return p, nil
}
package ping

import (
	"bytes"
	"encoding/binary"
	"fmt"
	network_tools "github.com/francisar/component/network/pkg/tools"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"time"
)

type pingV4 struct {
	Id int
	Seq int
	conn *net.IPConn
	raddr *net.IPAddr
	timeout time.Duration
	dataSize int
}

func NewPingV4(raddr *net.IPAddr, timeout time.Duration, Identify int) (Ping, error) {
	con, err := net.DialIP("ip4:icmp", nil, raddr)
	if err != nil {
		return nil, err
	}
	ping := pingV4{
		Id: Identify,
		conn: con,
		raddr: raddr,
		Seq: 0,
		timeout:timeout,
		dataSize:maxDataSize,
	}
	return &ping, nil
}

func (p *pingV4)SendIcmp(data []byte) error {
	dataSize := len(data)
	if dataSize > maxDataSize {
		msg := fmt.Sprintf("data size exceed the max data size limit:%d", maxDataSize)
		err := NewIcmpError(msg)
		return err
	}
	sendTime := network_tools.TimeToBytes(time.Now())
	t := append(sendTime, data...)
	body := &icmp.Echo{
		ID:   p.Id,
		Seq:  p.Seq,
		Data: t,
	}

	msg := &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: body,
	}
	msgBytes, err := msg.Marshal(nil)
	p.dataSize = len(msgBytes)
	if err != nil {
		icmpErr := NewIcmpError("data Marshal failed").withError(err)
		return icmpErr
	}
	var buffer bytes.Buffer
	err = binary.Write(&buffer, binary.BigEndian, msgBytes)
	if err != nil {
		icmpErr := NewIcmpError("data to binary buffer failed").withError(err)
		return icmpErr
	}
	_, err = p.conn.Write(buffer.Bytes())
	if err != nil {
		icmpErr := NewIcmpError("send data faild").withError(err)
		return icmpErr
	}
	p.Seq += 1
	return nil
}

func (p *pingV4)RecvIcmp() (*ICMPResponse,error) {
	err := p.conn.SetReadDeadline(time.Now().Add(p.timeout))
	if err != nil {
		icmpErr := NewIcmpError("set recv time out").withError(err)
		return nil, icmpErr
	}
	//构建接受的比特数组
	rec := make([]byte,maxDataSize)
	//读取连接返回的数据，将数据放入rec中
	n,err := p.conn.Read(rec)
	if err != nil {
		icmpErr := NewIcmpError("recv data faild").withError(err)
		return nil, icmpErr
	}
	ipv4header,err := icmp.ParseIPv4Header(rec)
	if err != nil {
		icmpErr := NewIcmpError("ParseIPv4Header faild").withError(err)
		return nil, icmpErr
	}
	echoMsg,err := icmp.ParseMessage(protocolICMP,rec[ipv4header.Len:n])
	icmpType, ok := echoMsg.Type.(ipv4.ICMPType)
	if !ok {
		icmpErr := NewIcmpError("resovle icmp type failed")
		return nil, icmpErr
	}
	responseData, err := echoMsg.Body.Marshal(int(icmpType))
	if err != nil {
		icmpErr := NewIcmpError("Marshal body failed").withError(err)
		return nil, icmpErr
	}
	switch icmpType {
	case ipv4.ICMPTypeEchoReply:
		ID := int(binary.BigEndian.Uint16(responseData[0:2]))
		Seq :=  int(binary.BigEndian.Uint16(responseData[2:4]))
		message := responseData[4:]
		timestamp := network_tools.BytesToTime(message[:timeSliceLength])
		rtt := time.Now().Sub(timestamp)
		icmpResponse := ICMPResponse{
			Code: echoMsg.Code,
			Rtt: rtt,
			Body: message[timeSliceLength:],
			Dst: ipv4header.Src,
			Src: ipv4header.Dst,
			TTL: ipv4header.TTL,
			ID: ID,
			Seq: Seq,

		}
		return &icmpResponse, nil
	default:
		msg := fmt.Sprintf("Un expected ICMP type:%s", icmpType.String())
		icmpErr := NewIcmpError(msg)
		return nil, icmpErr
	}

}
package ping

import (
	"bytes"
	"encoding/binary"
	"fmt"
	network_tools "github.com/francisar/component/network/pkg/tools"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv6"
	"net"
	"time"
)

type udpV6 struct {
	Id int
	Seq int
	conn *icmp.PacketConn
	raddr *net.IPAddr
	timeout time.Duration
	dataSize int
}


func NewUdpV6(raddr *net.IPAddr, timeout time.Duration, Identify int) (Ping, error) {
	con, err := icmp.ListenPacket("udp6", "0.0.0.0")
	if err != nil {
		return nil, err
	}

	err = con.IPv6PacketConn().SetControlMessage(ipv6.FlagHopLimit, true)
	if err != nil {
		return nil, err
	}
	ping := udpV6{
		Id: Identify,
		conn: con,
		raddr: raddr,
		Seq: 0,
		timeout:timeout,
		dataSize:maxDataSize,
	}
	return &ping, nil
}


func (p *udpV6)SendIcmp(data []byte) error {
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
		Type: ipv6.ICMPTypeEchoRequest,
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
	dst := &net.UDPAddr{IP: p.raddr.IP, Zone: p.raddr.Zone}
	_, err = p.conn.WriteTo(buffer.Bytes(), dst)
	if err != nil {
		icmpErr := NewIcmpError("send data faild").withError(err)
		return icmpErr
	}
	p.Seq += 1
	return nil
}

func (p *udpV6)RecvIcmp() (*ICMPResponse,error) {
	err := p.conn.SetReadDeadline(time.Now().Add(p.timeout))
	if err != nil {
		icmpErr := NewIcmpError("set recv time out failed").withError(err)
		return nil, icmpErr
	}
	//构建接受的比特数组
	rec := make([]byte,maxDataSize)
	n, cm, _, err := p.conn.IPv6PacketConn().ReadFrom(rec)
	if err != nil {
		icmpErr := NewIcmpError("recv data faild").withError(err)
		return nil, icmpErr
	}
	var sourceAddr net.IP
	var ttl int
	if cm != nil {
		sourceAddr = cm.Dst
		ttl = cm.HopLimit
	}
	echoMsg,err := icmp.ParseMessage(protocolIPv6ICMP,rec[:n])
	icmpType, ok := echoMsg.Type.(ipv6.ICMPType)
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
	case ipv6.ICMPTypeEchoReply:
		ID := int(binary.BigEndian.Uint16(responseData[0:2]))
		Seq :=  int(binary.BigEndian.Uint16(responseData[2:4]))
		message := responseData[4:]
		timestamp := network_tools.BytesToTime(message[:timeSliceLength])
		rtt := time.Now().Sub(timestamp)
		icmpResponse := ICMPResponse{
			Code: echoMsg.Code,
			Rtt: rtt,
			Body: message[timeSliceLength:],
			Dst: p.raddr.IP,
			Src: sourceAddr,
			TTL: ttl,
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

func (p *udpV6)Close() error {
	return p.conn.Close()
}
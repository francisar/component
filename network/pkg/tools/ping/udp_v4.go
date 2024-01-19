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

type udpV4 struct {
	Id int
	Seq int
	conn *icmp.PacketConn
	laddr *net.IPAddr
	raddr *net.IPAddr
	timeout time.Duration
	dataSize int
}


func NewUdpV4(laddr *net.IPAddr, raddr *net.IPAddr,  timeout time.Duration, Identify int) (Ping, error) {
	con, err := icmp.ListenPacket("udp4", laddr.IP.String())
	if err != nil {
		return nil, err
	}

	err = con.IPv4PacketConn().SetControlMessage(ipv4.FlagTTL, true)
	if err != nil {
		return nil, err
	}
	ping := udpV4{
		Id: Identify,
		conn: con,
		laddr: laddr,
		raddr: raddr,
		Seq: 0,
		timeout:timeout,
		dataSize:maxDataSize,
	}
	return &ping, nil
}


func (p *udpV4)SendIcmp(data []byte) error {
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
	dst := &net.UDPAddr{IP: p.raddr.IP}
	_, err = p.conn.WriteTo(buffer.Bytes(), dst)
	if err != nil {
		icmpErr := NewIcmpError("send data faild").withError(err)
		return icmpErr
	}
	p.Seq += 1
	return nil
}

func (p *udpV4)RecvIcmp() (*ICMPResponse,error) {
	err := p.conn.SetReadDeadline(time.Now().Add(p.timeout))
	if err != nil {
		icmpErr := NewIcmpError("set recv time out failed").withError(err)
		return nil, icmpErr
	}
	//构建接受的比特数组
	rec := make([]byte,maxDataSize)

	n, _, err := p.conn.ReadFrom(rec)
	// n,_,err := p.conn.ReadFrom(rec)
	if err != nil {
		if e, ok := err.(net.Error); ok && e.Timeout() {
			icmpResponse := ICMPResponse{
				Src: p.laddr.IP,
				Dst: p.raddr.IP,
				Code: -1,
				TTL: -1,
				ID: p.Id,
				Seq: p.Seq,
				Timeout: true,
			}
			return &icmpResponse, nil
		}

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
			Timeout: false,
			ICMPRet: icmpType,
		}
		return &icmpResponse, nil
	case ipv4.ICMPTypeDestinationUnreachable:
		message := responseData[4:]
		icmpResponse := ICMPResponse{
			ICMPRet: icmpType,
			Code: echoMsg.Code,
			Body: message[timeSliceLength:],
			Src: ipv4header.Dst,
			Timeout: false,
			TTL: ipv4header.TTL,
		}
		prePacketHeader,err := icmp.ParseIPv4Header(message)
		if err == nil {
			icmpResponse.Src = prePacketHeader.Src
			icmpResponse.Dst = prePacketHeader.Dst
			preEchoMsg,MsgErr := icmp.ParseMessage(protocolICMP,rec[prePacketHeader.Len:])
			if MsgErr == nil {
				echoByte, _ := preEchoMsg.Body.Marshal(ipv4.ICMPTypeEcho.Protocol())
				echo, tempErr := parseEcho(protocolICMP, ipv4.ICMPTypeEcho, echoByte)
				if tempErr == nil {
					icmpResponse.ID = echo.ID
					icmpResponse.Seq = echo.Seq
				}
			}
			icmpResponse.Body = message[prePacketHeader.Len:]
		}
		return &icmpResponse, nil
	default:
		msg := fmt.Sprintf("Un expected ICMP type:%s", icmpType.String())
		icmpErr := NewIcmpError(msg)
		return nil, icmpErr
	}
}

func (p *udpV4)Close() error {
	return p.conn.Close()
}

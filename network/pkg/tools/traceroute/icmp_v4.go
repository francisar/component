package traceroute

import (
	"bytes"
	"encoding/binary"
	"fmt"
	icmpbody "github.com/francisar/component/network/pkg/protocol/icmp"
	network_tools "github.com/francisar/component/network/pkg/tools"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"time"
)

type traceIcmpV4 struct {
	conn *net.PacketConn
	raddr *net.IPAddr
	timeout time.Duration
	id int
	seq int
}

func NewTraceIcmpV4(raddr *net.IPAddr, timeout time.Duration) (TraceRoute,error) {
	con, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, err
	}

	traceIcmpV4Impl := traceIcmpV4{
		conn: &con,
		raddr: raddr,
		timeout: timeout,
		id: network_tools.GetIdentify(),
		seq: network_tools.GetInitSeqId(),
	}
	return &traceIcmpV4Impl, nil
}


func (t *traceIcmpV4)SendPacket(ttl int) (*TraceResponse, error) {
	sendTime := time.Now()
	sendTimebyte := network_tools.TimeToBytes(sendTime)
	p := ipv4.NewPacketConn(*t.conn)
	if err := p.SetControlMessage(ipv4.FlagTTL|ipv4.FlagSrc|ipv4.FlagDst|ipv4.FlagInterface, true); err != nil {
		return nil, err
	}
	err := p.SetTTL(ttl)
	if err != nil {
		traceErr := NewTraceError("set ttl failed").withError(err)
		return nil,traceErr
	}
	body := &icmp.Echo{
		ID:   t.id,
		Seq:  t.seq,
		Data: sendTimebyte,
	}

	msg := &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: body,
	}
	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		traceErr := NewTraceError("data Marshal failed").withError(err)
		return nil,traceErr
	}
	var buffer bytes.Buffer
	err = binary.Write(&buffer, binary.BigEndian, msgBytes)
	if err != nil {
		traceErr := NewTraceError("data to binary buffer failed").withError(err)
		return nil,traceErr
	}
	cm := &ipv4.ControlMessage{
		TTL: ttl,
	}
	_, err = p.WriteTo(buffer.Bytes(), cm, t.raddr)
	if err != nil {
		traceErr := NewTraceError("send data faild").withError(err)
		return nil, traceErr
	}
	resp := TraceResponse{
		TTL: ttl,
		Dst: t.raddr.IP,
		ID: t.id,
		Seq: t.seq,
		SendTime:sendTime,
	}
	t.seq += 1
	return &resp, nil
}

func (t *traceIcmpV4)RecievePacket() (*TraceResponse, error) {
	rb := make([]byte, 1500)
	n, _, err := (*t.conn).ReadFrom(rb)
	if err != nil {
		traceErr := NewTraceError("read data failed").withError(err)
		return nil, traceErr
	}
	ipv4header,err := icmp.ParseIPv4Header(rb)
	if err != nil {
		traceErr := NewTraceError("ParseIPv4Header failed").withError(err)
		return nil, traceErr
	}
	rm, err := icmp.ParseMessage(icmpbody.ProtocolICMP, rb[ipv4header.Len:n])
	if err != nil {
		traceErr := NewTraceError("icmp ParseMessage failed").withError(err)
		return nil, traceErr
	}
	switch rm.Type {
	case ipv4.ICMPTypeTimeExceeded:
		body,err := rm.Body.Marshal(icmpbody.ProtocolICMP)
		if err != nil {
			traceErr := NewTraceError("icmp body Marshal failed").withError(err)
			return nil, traceErr
		}
		timeExceededBody,err := icmpbody.ParseTimeExceedeBody(body)
		if err != nil {
			traceErr := NewTraceError("icmp ParseTimeExceedeBody failed").withError(err)
			return nil, traceErr
		}
		resp := TraceResponse{
			TTL: timeExceededBody.IpHeader.TTL,
			Dst: timeExceededBody.IpHeader.Src,
			Src: timeExceededBody.IpHeader.Dst,
		}

		msg, err := icmp.ParseMessage(icmpbody.ProtocolICMP, timeExceededBody.Data)
		if err != nil {
			traceErr := NewTraceError("icmp ICMPTypeTimeExceeded data ParseMessage failed").withError(err)
			return nil, traceErr
		}
		if msg.Type == ipv4.ICMPTypeEcho {
			msgBody, err := msg.Body.Marshal(icmpbody.ProtocolICMP)
			if err != nil {
				traceErr := NewTraceError("icmp ICMPTypeTimeExceeded ICMPTypeEcho data Marshal failed").withError(err)
				return nil, traceErr
			}
			echo, err := icmpbody.ParseEchoBody(msgBody)
			if err != nil {
				traceErr := NewTraceError("icmp ICMPTypeTimeExceeded ICMPTypeEcho ParseEchoBody failed").withError(err)
				return nil, traceErr
			}
			resp.ID = echo.ID
			resp.Seq = echo.Seq
			return &resp, nil
		}

	case ipv4.ICMPTypeEchoReply:
		body,err := rm.Body.Marshal(icmpbody.ProtocolICMP)
		if err != nil {
			traceErr := NewTraceError("icmp ICMPTypeEchoReply body Marshal failed").withError(err)
			return nil, traceErr
		}
		echoBody,err := icmpbody.ParseEchoBody(body)
		if err != nil {
			traceErr := NewTraceError("icmp ICMPTypeEchoReply ParseEchoBody failed").withError(err)
			return nil, traceErr
		}
		resp := TraceResponse{
			TTL: ipv4header.TTL,
			Dst: ipv4header.Src,
			Src: ipv4header.Dst,
			Seq: echoBody.Seq,
			ID: echoBody.ID,
		}
		return &resp, nil
	default:
		msg := fmt.Sprintf("unexpected icmp type %v failed", rm.Type)
		traceErr := NewTraceError(msg)
		return nil, traceErr
	}
	traceErr := NewTraceError("unknown err")
	return nil, traceErr
}
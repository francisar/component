package traceroute

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"net"
)

type trace struct {
	traceStatistics map[string]*TraceResponse
	traceConn TraceRoute
}

func NewTrace(raddr *net.IPAddr) (*trace, error) {
	if raddr.IP.To4() != nil {
		traceConn,err := NewTraceIcmpV4(raddr, TimeOut)
		if err != nil {
			return nil, err
		}
		return &trace{
			traceConn: traceConn,
			traceStatistics: make(map[string]*TraceResponse, MaxTTL),
		},nil
	}
	return nil, NewTraceError("unKnown err")
}

func (t *trace)TraceSend() error {
	var errs *multierror.Error
	for i := 0; i< MaxTTL; i++ {
		resp, err := t.traceConn.SendPacket(i)
		if err != nil {
			errs = multierror.Append(err, errs)
			continue
		}
		index := fmt.Sprintf("%d_%d", resp.ID, resp.Seq)
		t.traceStatistics[index] = resp
	}
	if errs !=  nil && errs.Len() != 0 {
		return errs
	}
	return nil
}

func (t *trace)TraceRecieve() error  {
	var errs *multierror.Error
	for i := 0; i< MaxTTL; i++ {
		resp, err := t.traceConn.RecievePacket()
		if err != nil {
			errs = multierror.Append(err, errs)
			continue
		}
		index := fmt.Sprintf("%d_%d", resp.ID, resp.Seq)
		sendInfo, ok := t.traceStatistics[index]
		if !ok {
			fmt.Println(resp.Seq, resp.ID, resp.Src, resp.Dst, resp.TTL)
			traceErr := NewTraceError("recieve invalid Id and seq")
			errs = multierror.Append(traceErr, errs)
			continue
		}
		sendInfo.ReceiveTime = resp.ReceiveTime
	}
	if errs.Len() != 0 {
		return errs
	}
	return nil
}

func (t *trace)TraceResult() map[string]*TraceResponse {
	return t.traceStatistics
}
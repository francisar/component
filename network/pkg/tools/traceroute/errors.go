package traceroute

import "fmt"

type traceError struct {
	msg string
}

func (i *traceError)Error() string {
	return i.msg
}


func (i *traceError)withError(err error) *traceError {
	i.msg = fmt.Sprintf("%s, with error:%s", i.msg, err.Error())
	return i
}

func NewTraceError(msg string) *traceError {
	return &traceError{
		msg: msg,
	}
}
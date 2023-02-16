package ping

import "fmt"

type icmpError struct {
	msg string
}

func (i *icmpError)Error() string {
	return i.msg
}


func (i *icmpError)withError(err error) *icmpError {
	i.msg = fmt.Sprintf("%s, with error:%s", i.msg, err.Error())
	return i
}

func NewIcmpError(msg string) *icmpError {
	return &icmpError{
		msg: msg,
	}
}
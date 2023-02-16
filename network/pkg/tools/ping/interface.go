package ping

type Ping interface {
	SendIcmp(data []byte) error
	RecvIcmp() (*ICMPResponse,error)
}
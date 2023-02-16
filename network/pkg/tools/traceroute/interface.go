package traceroute




type TraceRoute interface {
	SendPacket(ttl int) (*TraceResponse, error)
	RecievePacket() (*TraceResponse, error)
}

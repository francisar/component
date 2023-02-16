package icmp

import "errors"

var (
	errInvalidProtocol  = errors.New("invalid protocol")
	errMessageTooShort  = errors.New("message too short")
	errHeaderTooShort   = errors.New("header too short")
	errBufferTooShort   = errors.New("buffer too short")
	errInvalidBody      = errors.New("invalid body")
)

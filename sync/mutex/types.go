package mutex

import "time"

type LockKey struct {
	Key string
	TimeOut time.Duration
	Meta string
}

type Status struct {
	Success bool
	Meta interface{}
	Err error
	ExpireAt time.Time
}


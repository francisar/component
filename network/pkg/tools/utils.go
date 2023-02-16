package tools

import (
	"math/rand"
	"os"
	"time"
)



func BytesToTime(b []byte) time.Time {
	var nsec int64
	for i := uint8(0); i < 8; i++ {
		nsec += int64(b[i]) << ((7 - i) * 8)
	}
	return time.Unix(nsec/1000000000, nsec%1000000000)
}


func TimeToBytes(t time.Time) []byte {
	nsec := t.UnixNano()
	b := make([]byte, 8)
	for i := uint8(0); i < 8; i++ {
		b[i] = byte((nsec >> ((7 - i) * 8)) & 0xff)
	}
	return b
}

func GetIdentify() int {
	return os.Getppid()
}

func GetInitSeqId() int {
	rand.Seed(time.Now().Unix())
	return rand.Int()
}

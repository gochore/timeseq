package timeseq

import (
	"encoding/binary"
	"time"
)

type Interval struct {
	NotBefore *time.Time
	NotAfter  *time.Time
}

type timeKey [16]byte

func (k timeKey) Get() time.Time {
	return time.Unix(int64(binary.BigEndian.Uint64(k[:8])), int64(binary.BigEndian.Uint64(k[8:])))
}

func newTimeKey(t time.Time) timeKey {
	var ret [16]byte
	binary.BigEndian.PutUint64(ret[:8], uint64(t.Unix()))
	binary.BigEndian.PutUint64(ret[8:], uint64(t.UnixNano()))
	return ret
}

package timeseq

import (
	"encoding/binary"
	"time"
)

type timeKey [12]byte

func (k timeKey) Time() time.Time {
	return time.Unix(int64(binary.BigEndian.Uint64(k[:8])), int64(binary.BigEndian.Uint32(k[8:])))
}

func newTimeKey(t time.Time) timeKey {
	var ret [12]byte
	binary.BigEndian.PutUint64(ret[:8], uint64(t.Unix()))
	binary.BigEndian.PutUint32(ret[8:], uint32(t.Nanosecond()))
	return ret
}

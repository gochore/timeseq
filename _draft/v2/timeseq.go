package timeseq

import (
	"encoding/binary"
	"sort"
	"time"
)

type Slice interface {
	// return length
	Len() int
	// swap items
	Swap(i, j int)
	// return time of item i
	Time(i int) time.Time
}

type sortableSlice struct {
	Slice
}

func (s sortableSlice) Less(i, j int) bool {
	return s.Time(i).Before(s.Time(j))
}

// Sort will sort seq by time
func Sort(slice Slice) {
	sort.Sort(sortableSlice{Slice: slice})
}

type timeKey [16]byte

func (k timeKey) Get() time.Time {
	return time.Unix(int64(binary.BigEndian.Uint64(k[:8])), int64(binary.BigEndian.Uint64(k[8:])))
}

func (k timeKey) Put(t time.Time) {
	binary.BigEndian.PutUint64(k[:8], uint64(t.Unix()))
	binary.BigEndian.PutUint64(k[8:], uint64(t.UnixNano()))
}

type Interval struct {
	NotBefore *time.Time
	NotAfter  *time.Time
}

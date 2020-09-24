package timeseq

import (
	"encoding/binary"
	"sort"
	"time"
)

type Interface interface {
	// return length
	Len() int
	// swap items
	Swap(i, j int)
	// return time of item i
	Time(i int) time.Time
	// return Slice[i:j]
	Slice(i, j int) Interface
}

type sortable struct {
	Interface
}

func (s sortable) Less(i, j int) bool {
	return s.Time(i).Before(s.Time(j))
}

// Sort will sort slice by time
func Sort(slice Interface) {
	sort.Stable(sortable{Interface: slice})
}

func IsSorted(slice Interface) bool {
	return sort.IsSorted(sortable{Interface: slice})
}

func Range(slice Interface, interval Interval) Interface {
	i := 0
	if interval.NotBefore != nil {
		i = sort.Search(slice.Len(), func(i int) bool {
			return !slice.Time(i).Before(*interval.NotBefore)
		})
	}
	j := slice.Len()
	if interval.NotAfter != nil {
		j = sort.Search(slice.Len(), func(j int) bool {
			return !slice.Time(j).Before(*interval.NotAfter)
		})
		if j < slice.Len() && slice.Time(j).Equal(*interval.NotAfter) {
			j++
		}
	}
	return slice.Slice(i, j)
}

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

package timeseq

import (
	"encoding/binary"
	"sort"
	"time"
)

// Interface is a type which can be sorted according to time
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

// Interval indicates a continuous time range
type Interval struct {
	NotBefore *time.Time
	NotAfter  *time.Time
}

// Contain return if time t is in the interval
func (i Interval) Contain(t time.Time) bool {
	if i.NotAfter != nil && t.After(*i.NotAfter) {
		return false
	}
	if i.NotBefore != nil && t.Before(*i.NotBefore) {
		return false
	}
	return true
}

// BeginAt is alias of AfterOrEqual
func (i Interval) BeginAt(t time.Time) Interval {
	return i.AfterOrEqual(t)
}

// EndAt is alias of BeforeOrEqual
func (i Interval) EndAt(t time.Time) Interval {
	return i.BeforeOrEqual(t)
}

// BeforeOrEqual return a new Interval which not before t
func (i Interval) BeforeOrEqual(t time.Time) Interval {
	return Interval{
		NotBefore: i.NotBefore,
		NotAfter:  &t,
	}
}

// AfterOrEqual return a new Interval which not after t
func (i Interval) AfterOrEqual(t time.Time) Interval {
	return Interval{
		NotBefore: &t,
		NotAfter:  i.NotAfter,
	}
}

// Before return a new Interval which before t
func (i Interval) Before(t time.Time) Interval {
	t = t.Add(-1)
	return Interval{
		NotBefore: i.NotBefore,
		NotAfter:  &t,
	}
}

// After return a new Interval which after t
func (i Interval) After(t time.Time) Interval {
	t = t.Add(1)
	return Interval{
		NotBefore: &t,
		NotAfter:  i.NotAfter,
	}
}

// BeginAt is alias of AfterOrEqual
func BeginAt(t time.Time) Interval {
	return AfterOrEqual(t)
}

// EndAt is alias of BeforeOrEqual
func EndAt(t time.Time) Interval {
	return BeforeOrEqual(t)
}

// BeforeOrEqual return a new Interval which not before t
func BeforeOrEqual(t time.Time) Interval {
	return Interval{
		NotAfter: &t,
	}
}

// AfterOrEqual return a new Interval which not after t
func AfterOrEqual(t time.Time) Interval {
	return Interval{
		NotBefore: &t,
	}
}

// Before return a new Interval which before t
func Before(t time.Time) Interval {
	t = t.Add(-1)
	return Interval{
		NotAfter: &t,
	}
}

// After return a new Interval which after t
func After(t time.Time) Interval {
	t = t.Add(1)
	return Interval{
		NotBefore: &t,
	}
}

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

package timeseq

import (
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

// Less implements sort.Interface.Less()
func (s sortable) Less(i, j int) bool {
	return s.Time(i).Before(s.Time(j))
}

// Sort will sort slice by time
func Sort(slice Interface) {
	sort.Stable(sortable{Interface: slice})
}

// IsSorted reports whether data is sorted.
func IsSorted(slice Interface) bool {
	return sort.IsSorted(sortable{Interface: slice})
}

// Range return a sub slice of given sorted slice according to the interval
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

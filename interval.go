package timeseq

import "time"

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

// String returns the interval formatted using the RFC3339 format string
func (i *Interval) String() string {
	return i.Format(time.RFC3339)
}

// Format returns a textual representation of the time value formatted according to layout
func (i *Interval) Format(layout string) string {
	notBefore, notAfter := "nil", "nil"
	if i.NotBefore != nil {
		notBefore = i.NotBefore.Format(layout)
	}
	if i.NotAfter != nil {
		notAfter = i.NotAfter.Format(layout)
	}
	return notBefore + "~" + notAfter
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

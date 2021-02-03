package timeseq

import "time"

// Interval indicates a continuous time range
type Interval struct {
	NotBefore *time.Time
	NotAfter  *time.Time
}

// Contain returns if time t is in the interval
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

// BeforeOrEqual returns a new Interval which not before t
func (i Interval) BeforeOrEqual(t time.Time) Interval {
	return Interval{
		NotBefore: i.NotBefore,
		NotAfter:  &t,
	}
}

// AfterOrEqual returns a new Interval which not after t
func (i Interval) AfterOrEqual(t time.Time) Interval {
	return Interval{
		NotBefore: &t,
		NotAfter:  i.NotAfter,
	}
}

// Before returns a new Interval which before t
func (i Interval) Before(t time.Time) Interval {
	t = t.Add(-1)
	return Interval{
		NotBefore: i.NotBefore,
		NotAfter:  &t,
	}
}

// After returns a new Interval which after t
func (i Interval) After(t time.Time) Interval {
	t = t.Add(1)
	return Interval{
		NotBefore: &t,
		NotAfter:  i.NotAfter,
	}
}

// Truncate returns the result of rounding interval down to a multiple of d (since the zero time).
func (i Interval) Truncate(d time.Duration) Interval {
	if i.NotBefore != nil {
		t := (*i.NotBefore).Truncate(d)
		if t.Before(*i.NotBefore) {
			t = t.Add(d)
		}
		i.NotBefore = &t
	}
	if i.NotAfter != nil {
		t := (*i.NotAfter).Truncate(d)
		i.NotAfter = &t
	}
	return i
}

// Duration returns the duration NotAfter - NotBefore,
// returns 0 if NotAfter is before or equal NotBefore,
// returns -1 if NotAfter or NotBefore if nil.
func (i Interval) Duration() time.Duration {
	if i.NotBefore == nil || i.NotAfter == nil {
		return -1
	}
	if !(*i.NotAfter).After(*i.NotBefore) {
		return 0
	}
	return (*i.NotAfter).Sub(*i.NotBefore)
}

// BeginAt is alias of AfterOrEqual
func BeginAt(t time.Time) Interval {
	return AfterOrEqual(t)
}

// EndAt is alias of BeforeOrEqual
func EndAt(t time.Time) Interval {
	return BeforeOrEqual(t)
}

// BeforeOrEqual returns a new Interval which not before t
func BeforeOrEqual(t time.Time) Interval {
	return Interval{
		NotAfter: &t,
	}
}

// AfterOrEqual returns a new Interval which not after t
func AfterOrEqual(t time.Time) Interval {
	return Interval{
		NotBefore: &t,
	}
}

// Before returns a new Interval which before t
func Before(t time.Time) Interval {
	t = t.Add(-1)
	return Interval{
		NotAfter: &t,
	}
}

// After returns a new Interval which after t
func After(t time.Time) Interval {
	t = t.Add(1)
	return Interval{
		NotBefore: &t,
	}
}

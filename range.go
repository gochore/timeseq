package timeseq

import "time"

// Range indicates a continuous time range
type Range struct {
	NotBefore *time.Time
	NotAfter  *time.Time
}

// Contain returns if time t is in the range
func (i Range) Contain(t time.Time) bool {
	if i.NotAfter != nil && t.After(*i.NotAfter) {
		return false
	}
	if i.NotBefore != nil && t.Before(*i.NotBefore) {
		return false
	}
	return true
}

// String returns the range formatted using the RFC3339 format string
func (i *Range) String() string {
	return i.Format(time.RFC3339)
}

// Format returns a textual representation of the time value formatted according to layout
func (i *Range) Format(layout string) string {
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
func (i Range) BeginAt(t time.Time) Range {
	return i.AfterOrEqual(t)
}

// EndAt is alias of BeforeOrEqual
func (i Range) EndAt(t time.Time) Range {
	return i.BeforeOrEqual(t)
}

// BeforeOrEqual returns a new Range which not before t
func (i Range) BeforeOrEqual(t time.Time) Range {
	return Range{
		NotBefore: i.NotBefore,
		NotAfter:  &t,
	}
}

// AfterOrEqual returns a new Range which not after t
func (i Range) AfterOrEqual(t time.Time) Range {
	return Range{
		NotBefore: &t,
		NotAfter:  i.NotAfter,
	}
}

// Before returns a new Range which before t
func (i Range) Before(t time.Time) Range {
	t = t.Add(-1)
	return Range{
		NotBefore: i.NotBefore,
		NotAfter:  &t,
	}
}

// After returns a new Range which after t
func (i Range) After(t time.Time) Range {
	t = t.Add(1)
	return Range{
		NotBefore: &t,
		NotAfter:  i.NotAfter,
	}
}

// Truncate returns the result of rounding range down to a multiple of d (since the zero time).
func (i Range) Truncate(d time.Duration) Range {
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
func (i Range) Duration() time.Duration {
	if i.NotBefore == nil || i.NotAfter == nil {
		return -1
	}
	if !(*i.NotAfter).After(*i.NotBefore) {
		return 0
	}
	return (*i.NotAfter).Sub(*i.NotBefore)
}

// BeginAt is alias of AfterOrEqual
func BeginAt(t time.Time) Range {
	return AfterOrEqual(t)
}

// EndAt is alias of BeforeOrEqual
func EndAt(t time.Time) Range {
	return BeforeOrEqual(t)
}

// BeforeOrEqual returns a new Range which not before t
func BeforeOrEqual(t time.Time) Range {
	return Range{
		NotAfter: &t,
	}
}

// AfterOrEqual returns a new Range which not after t
func AfterOrEqual(t time.Time) Range {
	return Range{
		NotBefore: &t,
	}
}

// Before returns a new Range which before t
func Before(t time.Time) Range {
	t = t.Add(-1)
	return Range{
		NotAfter: &t,
	}
}

// After returns a new Range which after t
func After(t time.Time) Range {
	t = t.Add(1)
	return Range{
		NotBefore: &t,
	}
}

package redTime

import "time"

type TimeImpl struct{}

func (TimeImpl) Now() time.Time {
	return time.Now()
}

var jst *time.Location

func (TimeImpl) JstLocation() *time.Location {
	if jst == nil {
		jst = time.FixedZone("Asia/Tokyo", 9*60*60)
	}
	return jst
}

package domain

import "time"

type Operation struct {
	ID          uint64
	UserID      uint64
	SegmentSlug string
	Operation   string
	Data        time.Time
}

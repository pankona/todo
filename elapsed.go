package kokizami

import (
	"time"
)

// Elapsed represents elapsed time of each Kizami
type Elapsed struct {
	Tag     string
	Desc    string
	Count   int
	Elapsed time.Duration
}

type SummaryRepository interface {
	ElapsedOfMonthByDesc(yyyymm string) ([]*Elapsed, error)
	ElapsedOfMonthByTag(yyyymm string) ([]*Elapsed, error)
}

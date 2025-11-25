package domain

import (
	"fmt"
)

// PixTrailer represents the trailer record for the PIX report
type PixTrailer struct {
	RecordType   string `fixed:"1,1" gorm:"column:record_type"`
	TotalRecords int64  `fixed:"2,11" gorm:"column:total_records"`
}

// NewPixTrailer creates a new PixTrailer instance
func NewPixTrailer(totalRecords int64) *PixTrailer {
	return &PixTrailer{
		RecordType:   "T",
		TotalRecords: totalRecords,
	}
}

// Format marshals the PixTrailer struct into a fixed-width format.
func (pt *PixTrailer) Format() string {
	ret := ""
	ret += fmt.Sprintf("%-1s", pt.RecordType)
	ret += fmt.Sprintf("%010d", pt.TotalRecords)
	return ret
}

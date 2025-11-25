package domain

import (
	"fmt"
	"time"
)

// PixHeader represents the header record for the PIX report
type PixHeader struct {
	RecordType    string    `fixed:"1,1"`
	DataTransacao time.Time `fixed:"2,11"`
}

// NewPixHeader creates a new PixHeader instance
func NewPixHeader(date time.Time) *PixHeader {
	return &PixHeader{
		RecordType:    "H",
		DataTransacao: date,
	}
}

// Format marshals the PixHeader struct into a fixed-width format.
func (ph *PixHeader) Format() string {
	ret := ""
	ret += fmt.Sprintf("%-1s", ph.RecordType)
	ret += ph.DataTransacao.Format("2006/01/02")
	return ret
}

package domain

import (
	"fmt"
	"time"

	"github.com/ianlopshire/go-fixedwidth"
)

type RankingHeader struct {
	FileName string `fixed:"1,8"`
	DateStr  string `fixed:"9,16"`
	Date     time.Time
	Acquirer string `fixed:"17,24"`
	Lines    int32  `fixed:"25,38"`
}

// Parse parses a line of text into a RankingHeader struct
func (rh *RankingHeader) Parse(line string) (*RankingHeader, error) {
	err := fixedwidth.Unmarshal([]byte(line), rh)
	if err != nil {
		return nil, err
	}
	rh.Date, err = time.Parse("20060102", rh.DateStr)
	if err != nil {
		return nil, err
	}
	return rh, nil
}

func (rh *RankingHeader) Validate(name string, lines int32) error {
	if rh.FileName != name {
		return fmt.Errorf("invalid file name: expected %s, got %s", name, rh.FileName)
	}
	if rh.Lines != lines {
		return fmt.Errorf("invalid line count: expected %d, got %d", lines, rh.Lines)
	}
	if rh.Acquirer != "47377613" {
		return fmt.Errorf("invalid acquirer: expected 47377613, got %s", rh.Acquirer)
	}
	return nil
}

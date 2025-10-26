package domain

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ianlopshire/go-fixedwidth"
	"github.com/lavinas/cadoc6334/internal/port"
)

// Segment represents a segment of a path.
type Segment struct {
	Name        string `fixed:"1,50" gorm:"column:nome_segmento"`
	Description string `fixed:"51,300" gorm:"column:descricao_segmento"`
	CodeStr     string `fixed:"301,303"`
	Code        int64  `gorm:"column:codigo_segmento"`
}

// NewSegment creates a new Segment instance
func NewSegment() *Segment {
	return &Segment{}
}

// Validate validates the Segment information.
func (s *Segment) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("invalid name in segment")
	}
	if s.Description == "" {
		return fmt.Errorf("invalid description in segment")
	}
	if s.Code <= 0 {
		return fmt.Errorf("invalid code in segment")
	}
	return nil
}

// TableName returns the table name for the Segment struct
func (s *Segment) TableName() string {
	return "cadoc_6334_segmentos"
}

// GetKey generates a unique key for the Segment record.
func (s *Segment) GetKey() string {
	return fmt.Sprintf("%03d", s.Code)
}

// FindAll retrieves all Segment records.
func (s *Segment) GetDB(repo port.Repository) (map[string]port.Report, error) {
	var records []*Segment
	err := repo.FindAll(&records)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]port.Report)
	for _, r := range records {
		ret[r.GetKey()] = r
	}
	return ret, nil
}

// Parse parses a fixed-width string into a Segment struct
func (s *Segment) Parse(line string) error {
	line, err := RemoveAccents(line)
	if err != nil {
		return fmt.Errorf("error removing accents: %w", err)
	}
	err = fixedwidth.Unmarshal([]byte(line), s)
	if err != nil {
		return err
	}
	if s.CodeStr == "" {
		s.Code = 0
		return nil
	}
	_, err = fmt.Sscanf(s.CodeStr, "%d", &s.Code)
	if err != nil {
		return fmt.Errorf("error parsing Code: %w", err)
	}
	return nil
}

// String returns a string representation of the Segment struct
func (s *Segment) String() string {
	return fmt.Sprintf("Name: %s, Description: %s, Code: %d", s.Name, s.Description, s.Code)
}

// ParseSegmentFile parses a file of segments into a slice of Segment structs
func (s *Segment) ParseSegmentFile(filename string) ([]*Segment, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	// read header
	if !scanner.Scan() {
		return nil, fmt.Errorf("file is empty")
	}
	headerLine := scanner.Text()
	header := &RankingHeader{}
	_, err = header.Parse(headerLine)
	if err != nil {
		return nil, fmt.Errorf("error parsing header: %w", err)
	}
	// read records
	segments := []*Segment{}
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		segment := &Segment{}
		err := segment.Parse(line)
		if err != nil {
			return nil, fmt.Errorf("error parsing line: %w", err)
		}
		segments = append(segments, segment)
		count++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if err := header.Validate("SEGMENTO", int64(count)); err != nil {
		return nil, err
	}
	return segments, nil
}

// GetParsedFile retrieves and maps Segment records from a file.
func (s *Segment) GetParsedFile(filename string) (map[string]port.Report, error) {
	fileSegments, err := s.ParseSegmentFile(filename)
	if err != nil {
		return nil, err
	}
	segmentMap := make(map[string]port.Report)
	for _, seg := range fileSegments {
		segmentMap[seg.GetKey()] = seg
	}
	return segmentMap, nil
}

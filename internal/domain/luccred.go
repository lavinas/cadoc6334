package domain

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ianlopshire/go-fixedwidth"
	"github.com/lavinas/cadoc6334/internal/port"
	"golang.org/x/text/encoding/charmap"
)

type LucrCred struct {
	Year               int64   `fixed:"1,4" gorm:"column:ano"`
	Quarter            int64   `fixed:"5,5" gorm:"column:trimestre"`
	DiscountRevenue    float64 `gorm:"column:receitataxadescontobruta;type:numeric(18,2)"`
	DiscountRevenueInt int64   `fixed:"6,17"`
	RentRevenue        float64 `gorm:"column:receitaaluguelequipamentosconectividade;type:numeric(18,2)"`
	RentRevenueInt     int64   `fixed:"18,29"`
	OtherRevenue       float64 `gorm:"column:receitaoutras"`
	OtherRevenueInt    int64   `fixed:"30,41"`
	InterchangeCost    float64 `gorm:"column:custotarifaintercambio"`
	InterchangeCostInt int64   `fixed:"42,53"`
	MarketingCost      float64 `gorm:"column:customarketingpropaganda"`
	MarketingCostInt   int64   `fixed:"54,65"`
	BrandAccessCost    float64 `gorm:"column:custotaxasacessobandeiras"`
	BrandAccessCostInt int64   `fixed:"66,77"`
	RiskCost           float64 `gorm:"column:custorisco"`
	RiskCostInt        int64   `fixed:"78,89"`
	ProcessingCost     float64 `gorm:"column:custoprocessamento"`
	ProcessingCostInt  int64   `fixed:"90,101"`
	OtherCost          float64 `gorm:"column:custooutros"`
	OtherCostInt       int64   `fixed:"102,113"`
}

// NewLucrCred creates a new LucrCred instance
func NewLucrCred() *LucrCred {
	return &LucrCred{}
}

// GetName gets name of the report
func (l *LucrCred) GetName() string {
	return "LUCRCRED"
}

// Format marshals the LucrCred struct into a fixed-width format.
func (l *LucrCred) Format() string {
	ret := ""
	ret += fmt.Sprintf("%04d", l.Year)
	ret += fmt.Sprintf("%01d", l.Quarter)
	// Convert float64 fields to int representation
	l.DiscountRevenueInt = int64(l.DiscountRevenue*100 + 0.5)
	l.RentRevenueInt = int64(l.RentRevenue*100 + 0.5)
	l.OtherRevenueInt = int64(l.OtherRevenue*100 + 0.5)
	l.InterchangeCostInt = int64(l.InterchangeCost*100 + 0.5)
	l.MarketingCostInt = int64(l.MarketingCost*100 + 0.5)
	l.BrandAccessCostInt = int64(l.BrandAccessCost*100 + 0.5)
	l.RiskCostInt = int64(l.RiskCost*100 + 0.5)
	l.ProcessingCostInt = int64(l.ProcessingCost*100 + 0.5)
	l.OtherCostInt = int64(l.OtherCost*100 + 0.5)
	ret += fmt.Sprintf("%012d", l.DiscountRevenueInt)
	ret += fmt.Sprintf("%012d", l.RentRevenueInt)
	ret += fmt.Sprintf("%012d", l.OtherRevenueInt)
	ret += fmt.Sprintf("%012d", l.InterchangeCostInt)
	ret += fmt.Sprintf("%012d", l.MarketingCostInt)
	ret += fmt.Sprintf("%012d", l.BrandAccessCostInt)
	ret += fmt.Sprintf("%012d", l.RiskCostInt)
	ret += fmt.Sprintf("%012d", l.ProcessingCostInt)
	ret += fmt.Sprintf("%012d", l.OtherCostInt)
	return ret
}

// Validate validates the LucrCred information.
func (l *LucrCred) Validate() error {
	if l.Year <= 0 {
		return fmt.Errorf("invalid year in LucrCred")
	}
	if l.Quarter <= 0 {
		return fmt.Errorf("invalid quarter in LucrCred")
	}
	if l.DiscountRevenue <= 0 {
		return fmt.Errorf("invalid discount revenue in LucrCred")
	}
	if l.RentRevenue <= 0 {
		return fmt.Errorf("invalid rent revenue in LucrCred")
	}
	if l.OtherRevenue <= 0 {
		return fmt.Errorf("invalid other revenue in LucrCred")
	}
	if l.InterchangeCost <= 0 {
		return fmt.Errorf("invalid interchange cost in LucrCred")
	}
	if l.MarketingCost <= 0 {
		return fmt.Errorf("invalid marketing cost in LucrCred")
	}
	if l.BrandAccessCost <= 0 {
		return fmt.Errorf("invalid brand access cost in LucrCred")
	}
	if l.RiskCost < 0 {
		return fmt.Errorf("invalid risk cost in LucrCred")
	}
	if l.ProcessingCost <= 0 {
		return fmt.Errorf("invalid processing cost in LucrCred")
	}
	if l.OtherCost <= 0 {
		return fmt.Errorf("invalid other cost in LucrCred")
	}
	return nil
}

// TableName returns the table name for the LucrCred struct
func (l *LucrCred) TableName() string {
	return "cadoc_6334_luccred"
}

// GetKey generates a unique key for the LucrCred record.
func (l *LucrCred) GetKey() string {
	return fmt.Sprintf("%04d%01d", l.Year, l.Quarter)
}

// GetDB retrieves all LucrCred records.
func (l *LucrCred) GetDB(repo port.Repository) (map[string]port.Report, error) {
	var records []*LucrCred
	err := repo.FindAll(&records, 0, 0, "")
	if err != nil {
		return nil, err
	}
	ret := make(map[string]port.Report)
	for _, r := range records {
		ret[r.GetKey()] = r
	}
	return ret, nil
}

// Parse parses a fixed-width string into a LucrCred struct
func (l *LucrCred) Parse(line string) error {
	err := fixedwidth.Unmarshal([]byte(line), l)
	if err != nil {
		return err
	}
	l.DiscountRevenue = float64(l.DiscountRevenueInt) / float64(100)
	l.RentRevenue = float64(l.RentRevenueInt) / float64(100)
	l.OtherRevenue = float64(l.OtherRevenueInt) / float64(100)
	l.InterchangeCost = float64(l.InterchangeCostInt) / float64(100)
	l.MarketingCost = float64(l.MarketingCostInt) / float64(100)
	l.BrandAccessCost = float64(l.BrandAccessCostInt) / float64(100)
	l.RiskCost = float64(l.RiskCostInt) / float64(100)
	l.ProcessingCost = float64(l.ProcessingCostInt) / float64(100)
	l.OtherCost = float64(l.OtherCostInt) / float64(100)
	return nil
}

// String returns a string representation of the LucrCred struct
func (l *LucrCred) String() string {
	return fmt.Sprintf("Year: %d, Quarter: %d, DiscountRevenue: %.2f, RentRevenue: %.2f, OtherRevenue: %.2f, InterchangeCost: %.2f, MarketingCost: %.2f, BrandAccessCost: %.2f, RiskCost: %.2f, ProcessingCost: %.2f, OtherCost: %.2f",
		l.Year, l.Quarter, l.DiscountRevenue, l.RentRevenue, l.OtherRevenue, l.InterchangeCost, l.MarketingCost, l.BrandAccessCost, l.RiskCost, l.ProcessingCost, l.OtherCost)
}

// ParseLucrCredFile parses the LucrCred.TXT file and returns a slice of LucrCred records.
func (l *LucrCred) ParseLucrCredFile(filePath string) ([]*LucrCred, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []*LucrCred
	decoder := charmap.ISO8859_1.NewDecoder()
	decodedReader := decoder.Reader(file)
	scanner := bufio.NewScanner(decodedReader)
	// Read and parse header
	if !scanner.Scan() {
		return nil, fmt.Errorf("file is empty")
	}
	headerLine := scanner.Text()
	header := &RankingHeader{}
	_, err = header.Parse(headerLine)
	if err != nil {
		return nil, fmt.Errorf("error parsing header: %w", err)
	}
	// Read and parse records
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		LucrCred := NewLucrCred()
		err := LucrCred.Parse(line)
		if err != nil {
			return nil, fmt.Errorf("error parsing line: %w", err)
		}
		records = append(records, LucrCred)
		count++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	// Validate header
	if err := header.Validate("LUCRCRED", int64(count)); err != nil {
		return nil, err
	}
	return records, nil
}

// GetParsedFile retrieves and parses the LucrCred.TXT file.
func (l *LucrCred) GetParsedFile(filename string) (map[string]port.Report, error) {
	records, err := l.ParseLucrCredFile(filename)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]port.Report)
	for _, r := range records {
		ret[r.GetKey()] = r
	}
	return ret, nil
}

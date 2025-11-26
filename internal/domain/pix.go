package domain

import (
	"fmt"
	"github.com/shopspring/decimal"
	"time"

	"github.com/ianlopshire/go-fixedwidth"
	"github.com/lavinas/cadoc6334/internal/port"
)

// Pix represents the PIX record structure.
type Pix struct {
	RecordType         string          `fixed:"1,1" gorm:"column:recordtype"`
	CodigoCliente      string          `fixed:"2,16" gorm:"column:codigocliente"`
	DataMovimento      time.Time       `fixed:"17,26" gorm:"column:data_movimento"`
	DataTransacao      time.Time       `fixed:"27,36" gorm:"column:datatransacao"`
	DataProcessamento  time.Time       `fixed:"37,46" gorm:"column:dataprocessamento"`
	CodigoBandeira     string          `fixed:"47,49" gorm:"column:codigobandeira"`
	CodigoProduto      string          `fixed:"50,51" gorm:"column:produto"`
	TipoParcelamento   string          `fixed:"52,52" gorm:"column:tp_parcela"`
	TipoTransacao      string          `fixed:"53,53" gorm:"column:tp_origem"`
	PlanoPagamento     string          `fixed:"54,55" gorm:"column:nu_parcela"`
	ValorBrutoOriginal decimal.Decimal `fixed:"56,72" gorm:"column:nu_valor"`
	TaxaMDROriginal    string          `fixed:"73,77" gorm:"column:nu_porc_transacao"`
	ValorMDROriginal   decimal.Decimal `fixed:"78,94" gorm:"column:valor_mdr"`
	TipoTecnologia     string          `fixed:"95,96" gorm:"column:tipo_tecnologia"`
	NumeroTerminal     string          `fixed:"97,104" gorm:"column:terminal_id"`
	CodigoAutorizacao  string          `fixed:"105,110" gorm:"column:auth_target"`
	NSU                string          `fixed:"111,130" gorm:"column:nsu_target"`
	NumeroECFPADQ      string          `fixed:"131,145" gorm:"column:numero_ec_fp_adq"`
	ArranjoPagamentoFP string          `fixed:"146,148" gorm:"column:empty_field"`
	CodigoFormaEntrada string          `fixed:"149,150" gorm:"column:forma_entrada"`
	Hora               time.Time       `fixed:"151,158" gorm:"column:hora_transacao;type:time"`
	Duplicated         bool            `gorm:"column:duplicated"`
}

// NewPix creates a new Pix instance
func NewPix() *Pix {
	return &Pix{}
}

// GetName gets name of the report
func (p *Pix) GetName() string {
	return "PIX"
}

// Format marshals the Pix struct into a fixed-width format.
func (p *Pix) Format() string {
	ret := ""
	ret += fmt.Sprintf("%-1s", p.RecordType)
	ret += fmt.Sprintf("%-15s", p.CodigoCliente)
	ret += p.DataMovimento.Format("2006/01/02")
	ret += p.DataTransacao.Format("2006/01/02")
	ret += p.DataProcessamento.Format("2006/01/02")
	ret += fmt.Sprintf("%-3s", p.CodigoBandeira)
	ret += fmt.Sprintf("%-2s", p.CodigoProduto)
	ret += fmt.Sprintf("%-1s", p.TipoParcelamento)
	ret += fmt.Sprintf("%-1s", p.TipoTransacao)
	ret += fmt.Sprintf("%-2s", p.PlanoPagamento)
	vbint := p.ValorBrutoOriginal.Mul(decimal.NewFromInt(100)).IntPart()
	ret += fmt.Sprintf("%017d", vbint)
	ret += fmt.Sprintf("%05s", p.TaxaMDROriginal)
	vmdr := p.ValorMDROriginal.Mul(decimal.NewFromInt(100)).IntPart()
	ret += fmt.Sprintf("%017d", vmdr)
	ret += fmt.Sprintf("%-2s", p.TipoTecnologia)
	ret += fmt.Sprintf("%-8s", p.NumeroTerminal)
	ret += fmt.Sprintf("%-6s", p.CodigoAutorizacao[:6])
	ret += fmt.Sprintf("%-20s", p.NSU[:20])
	ret += fmt.Sprintf("%-15s", p.NumeroECFPADQ)
	ret += fmt.Sprintf("%-3s", p.ArranjoPagamentoFP)
	ret += fmt.Sprintf("%-2s", p.CodigoFormaEntrada)
	ret += p.Hora.Format("15:04:05")
	return ret
}

// Validate validates the Pix information.
func (p *Pix) Validate() error {
	if p.RecordType == "" {
		return fmt.Errorf("invalid record type in pix")
	}
	if p.CodigoCliente == "" {
		return fmt.Errorf("invalid codigo cliente in pix")
	}
	return nil
}

// TableName returns the table name for the Pix struct.
func (p *Pix) TableName() string {
	return "pix_dimp"
}

// GetKey returns the unique key for the Pix record.
func (p *Pix) GetKey() string {
	return fmt.Sprintf("%s|%s", p.DataTransacao.Format("2006-01-02"), p.NSU)
}

// GetDB returns the database connection.
func (p *Pix) GetDB(repo port.Repository) (map[string]port.Report, error) {
	var records []*Pix
	err := repo.FindAll(&records, 0, 0, "datatransacao")
	if err != nil {
		return nil, err
	}
	result := make(map[string]port.Report)
	for _, record := range records {
		if record.Duplicated {
			continue
		}
		result[record.GetKey()] = record
	}
	return result, nil
}

// GetDB returns the database connection.
func (p *Pix) GetDBOrdered(repo port.Repository) ([]port.Report, error) {
	var records []*Pix
	err := repo.FindAll(&records, 0, 0, "datatransacao")
	if err != nil {
		return nil, err
	}
	var result []port.Report
	for _, record := range records {
		if record.Duplicated {
			continue
		}
		result = append(result, record)
	}
	return result, nil
}

// Parse parses the Pix data from a fixed-width file.
func (p *Pix) Parse(line string) error {
	err := fixedwidth.Unmarshal([]byte(line), p)
	if err != nil {
		return err
	}
	return nil
}

// String returns a string representation of the Pix struct.
func (p *Pix) String() string {
	return fmt.Sprintf("Pix{RecordType: %s, CodigoCliente: %s, DataMovimento: %s, DataTransacao: %s, DataProcessamento: %s, CodigoBandeira: %s, CodigoProduto: %s, TipoParcelamento: %s, TipoTransacao: %s, PlanoPagamento: %s, ValorBrutoOriginal: %s, TaxaMDROriginal: %s, ValorMDROriginal: %s, TipoTecnologia: %s, NumeroTerminal: %s, CodigoAutorizacao: %s, NSU: %s, NumeroECFPADQ: %s, ArranjoPagamentoFP: %s, CodigoFormaEntrada: %s, Hora: %s}",
		p.RecordType, p.CodigoCliente, p.DataMovimento.Format("2006-01-02"), p.DataTransacao.Format("2006-01-02"), p.DataProcessamento.Format("2006-01-02"), p.CodigoBandeira, p.CodigoProduto, p.TipoParcelamento, p.TipoTransacao, p.PlanoPagamento, p.ValorBrutoOriginal.String(), p.TaxaMDROriginal, p.ValorMDROriginal.String(), p.TipoTecnologia, p.NumeroTerminal, p.CodigoAutorizacao, p.NSU, p.NumeroECFPADQ, p.ArranjoPagamentoFP, p.CodigoFormaEntrada, p.Hora.Format("15:04:05"))
}

// ParsePixFile parses the PIX.TXT file and returns a slice of Pix records.
func (p *Pix) ParsePixFile(lines []string) ([]*Pix, error) {
	var records []*Pix
	for _, line := range lines {
		var pix Pix
		err := pix.Parse(line)
		if err != nil {
			return nil, err
		}
		records = append(records, &pix)
	}
	return records, nil
}

// GetParsePixFile parses the PIX.TXT file and returns a slice of Pix records.
func (p *Pix) GetParsedFile(filename string) (map[string]port.Report, error) {
	return nil, nil
}

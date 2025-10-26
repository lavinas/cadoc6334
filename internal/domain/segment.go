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
	Name        string `fixed:"1,50"`
	Description string `fixed:"51,300"`
	Code        int32  
	CodeStr	    string `fixed:"301,303"`
}

// NewSegment creates a new Segment instance
func NewSegment() *Segment {
	return &Segment{}
}

// TableName returns the table name for the Segment struct
func (s *Segment) TableName() string {
	return "cadoc_6334_segmento"
}

// GetKey generates a unique key for the Segment record.
func (s *Segment) GetKey() string {
	return fmt.Sprintf("%03d", s.Code)
}

// FindAll retrieves all Segment records.
func (s *Segment) FindAll(repo port.Repository) (map[string]port.Report, error) {
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
	fmt.Println("Parsing line:", line)
	err := fixedwidth.Unmarshal([]byte(line), s)
	if err != nil {
		return err
	}
	// Convert CodeStr to Code
	_, err = fmt.Sscanf(s.CodeStr, "%03d", &s.Code)
	if err != nil {
		return fmt.Errorf("error parsing Code: %w", err)
	}
	fmt.Printf("Parsed segment: %+v\n", s)
	return nil
}

// String returns a string representation of the Segment struct
func (s *Segment) String() string {
	return fmt.Sprintf("Name: %s, Description: %s, Code: %d", s.Name, s.Description, s.Code)
}

// LoadSegments loads segments from fixed string
func (s *Segment) LoadSegments() ([]*Segment, error) {
	ret := []*Segment{}
	ret = append(ret, &Segment{
		Name:        "Cuidados pessoais",
		Description: "Loja de cosméticos;Navalha elétrica - venda e serviços",
		Code:        401,
	})
	ret = append(ret, &Segment{
		Name:        "Bares e Restaurantes",
		Description: "Bares, pubs e casas noturnas;Bares de sinuca",
		Code:        402,
	})
	ret = append(ret, &Segment{
		Name:        "Companhias aéreas e afins",
		Description: "Aeroportos e serviços ligados a aeronaves",
		Code:        403,
	})
	ret = append(ret, &Segment{
		Name:        "Cultura e Esportes",
		Description: "Cinemas, produções cinematográficas;Academias / clubes",
		Code:        404,
	})
	ret = append(ret, &Segment{
		Name:        "Educação",
		Description: "Universidades e faculdades;Escola de negócios/vocações;Colégios",
		Code:        405,
	})
	ret = append(ret, &Segment{
		Name:        "Eletrônicos e eletrodomésticos",
		Description: "Computadores, equipamentos e softwares;Produtos digitais - aplicativos de software (exceto jogos);Lojas de eletrodomésticos",
		Code:        406,
	})
	ret = append(ret, &Segment{
		Name:        "Farmácias e Cuidados com a saúde",
		Description: "Aparelhos auditivos - vendas e serviços;Farmácias",
		Code:        407,
	})
	ret = append(ret, &Segment{
		Name:        "Grandes Atacadistas",
		Description: "Atacados de bebidas alcoólicas",
		Code:        408,
	})
	ret = append(ret, &Segment{
		Name:        "Outros Serviços e Profissionais Liberais",
		Description: "Corretores de imóveis;Serviço funerário;Consultoria empresarial e serviços de relações públicas;Outros serviços profissionais de especializados",
		Code:        409,
	})
	ret = append(ret, &Segment{
		Name:        "Jogos e Loteria",
		Description: "Corrida de cavalos licenciado;Cassinos, loterias e jogos de azar",
		Code:        410,
	})
	ret = append(ret, &Segment{
		Name:        "Livrarias e afins",
		Description: "Banca de jornal e provedor de notícias;Artigos de papelaria e suprimentos para escritório",
		Code:        411,
	})
	ret = append(ret, &Segment{
		Name:        "Alimentação",
		Description: "Loja de doces",
		Code:        412,
	})
	ret = append(ret, &Segment{
		Name:        "Móveis e construção",
		Description: "Demais serviços de reforma e construção;Piscinas e banheiras - serviços, suprimentos e vendas",
		Code:        413,
	})
	ret = append(ret, &Segment{
		Name:        "Pequenos supermercados e afins",
		Description: "Lojas especializadas não listadas anteriormente;Lojas de variedades",
		Code:        414,
	})
	ret = append(ret, &Segment{
		Name:        "Combustíveis e afins",
		Description: "Postos de gasolina;Revendedores de combustíveis",
		Code:        415,
	})
	ret = append(ret, &Segment{
		Name:        "Roupas, sapatos, acessórios e afins",
		Description: "Conserto de relógios e joias;Aluguel de roupas - fantasias, uniformes e roupas sociais",
		Code:        416,
	})
	ret = append(ret, &Segment{
		Name:        "Comércio e serviços em geral",
		Description: "Opticians, optical goods, and eyeglasses;Loja de moedas e selos",
		Code:        421,
	})
	ret = append(ret, &Segment{
		Name:        "Serviços Financeiros",
		Description: "Instituição financeira - agências e serviços;Corretores de residências móveis",
		Code:        422,
	})
	ret = append(ret, &Segment{
		Name:        "Outros",
		Description: "Armazenamento agrícola, refrigeração, bens domésticos;Telegrafo",
		Code:        423,
	})
	ret = append(ret, &Segment{
		Name:        "Instituições Financeiras",
		Description: "Bancos / lojas de poupança e inst. Financeira;Instituição financeira - caixa eletrônico",
		Code:        424,
	})
	ret = append(ret, &Segment{
		Name:        "Serviços Públicos",
		Description: "Multas (fines);Serviços governamentais",
		Code:        425,
	})
	ret = append(ret, &Segment{
		Name:        "Seguros",
		Description: "Marketing direto de seguros",
		Code:        426,
	})
	ret = append(ret, &Segment{
		Name:        "Utilities (inclui telecom)",
		Description: "Telefones e equipamentos de telecom.;Catálogo de varejo",
		Code:        427,
	})
	ret = append(ret, &Segment{
		Name:        "Subadquirentes",
		Description: "Catálogo de varejo",
		Code:        428,
	})
	return ret, nil
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
	if err := header.Validate("SEGMENTO", int32(count)); err != nil {
		return nil, err
	}
	return segments, nil
}

// GetLoaded retrieves and maps Segment records from mounted data.
func (s *Segment) GetLoaded() (map[string]port.Report, error) {
	loadedSegments, err := s.LoadSegments()
	if err != nil {
		return nil, err
	}
	if loadedSegments == nil {
		return nil, fmt.Errorf("no segment data loaded")
	}
	mapSegments := make(map[string]port.Report)
	for _, seg := range loadedSegments {
		mapSegments[seg.GetKey()] = seg
	}
	return mapSegments, nil
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


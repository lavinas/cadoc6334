package domain

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ianlopshire/go-fixedwidth"
)

// Segment represents a segment of a path.
type Segment struct {
	Name        string `fixed:"1,50"`
	Description string `fixed:"51,300"`
	Code        int32  
	CodeStr	    string `fixed:"301,303"`
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
func LoadSegments() ([]*Segment, error) {
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
func ParseSegmentFile(filename string) ([]*Segment, error) {
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

// ReconciliateSegments reconciles two slices of Segment structs and returns the differences
func ReconciliateSegments(filename string) {
	// load segments from fixed string
	fixedSegments, err := LoadSegments()
	if err != nil {
		fmt.Printf("Error loading fixed segments: %v\n", err)
		return
	}
	// parse segments from file
	parsedSegments, err := ParseSegmentFile(filename)
	if err != nil {
		fmt.Printf("Error parsing segment file: %v\n", err)
		return
	}
	if len(fixedSegments) != len(parsedSegments) {
		fmt.Printf("Segment count mismatch: fixed=%d, parsed=%d\n", len(fixedSegments), len(parsedSegments))
		return
	}
	map1 := map[int32]*Segment{}
	map2 := map[int32]*Segment{}
	for _, s := range fixedSegments {
		map1[s.Code] = s
	}
	for _, s := range parsedSegments {
		map2[s.Code] = s
	}
	for code, seg1 := range map1 {
		seg2, ok := map2[code]
		if !ok {
			continue
		}
		if seg1.String() != seg2.String() {
			fmt.Printf("Segment code %d mismatch:\n  fixed:  %s\n  parsed: %s\n", code, seg1.String(), seg2.String())
		}
	}
}

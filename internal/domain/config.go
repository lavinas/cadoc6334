package domain

// distributions
var (
	// period
	years    = []int32{2025}
	quarters = []int32{3}
	// functions
	funcValues = []string{"D", "C"}
	funcProp   = []float32{0.7, 0.3}
	// Brands
	brandValues = []int32{1, 2, 8}
	brandProp   = []float32{0.5, 0.3, 0.2}
	// capture
	captValues = []int32{2, 5}
	captProp   = []float32{0.7, 0.3}
	// credit installments
	instCredValues = []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	instCredProp   = []float32{0.22, 0.14, 0.13, 0.1, 0.08, 0.07, 0.05, 0.05, 0.04, 0.04, 0.04, 0.04}
	instCredFee    = []float32{0.00, 1.10, 1.10, 1.10, 1.10, 1.10, 2.15, 2.15, 2.15, 2.15, 2.15, 2.15}
	// debit installments
	instDebValues = []int32{1}
	instDebProp   = []float32{1}
	instDebFee    = []float32{-0.5}
	// segments
	segValues = []int32{401, 402, 403, 404, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 421, 422, 423, 424, 425, 426, 427, 428}
	segProp   = []float32{0.05, 0.05, 0.05, 0.05, 0.05, 0.05, 0.05, 0.05, 0.05, 0.05, 0.05, 0.05, 0.05, 0.05, 0.05, 0.05, 0.02, 0.02, 0.02, 0.02, 0.02, 0.02, 0.04, 0.04}
	// QttyRange
	avgTicket float32 = 150
	// ufs
	ufValues = []string{"SP", "MG", "RJ", "BA", "PR", "RS", "PE", "CE", "PA", "SC", "GO", "MA", "AM", "PB", "ES", "MT", "RN", "PI", "AL", "DF", "MS", "SE", "RO", "TO", "AC", "AP", "RR"}
	ufProp   = []float32{0.2159, 0.1002, 0.0807, 0.0697, 0.0557, 0.0526, 0.0448, 0.0434, 0.0408, 0.0384, 0.0348, 0.0329, 0.0202, 0.0195, 0.0193, 0.0182, 0.0162, 0.0158, 0.0151, 0.0140, 0.0137, 0.0108, 0.0082, 0.0074, 0.0041, 0.0039, 0.0030}
	// cardType
	cardtypeValues = []string{"P", "H", "C"}
	cardtypeProp   = []float32{0.8, 0.1, 0.1}

	//product
	prodValues = []int32{32, 33, 34, 35, 36, 37}
	prodProp   = []float32{0.5, 0.3, 0.1, 0.05, 0.025, 0.0025}
)

// ranking
var (
	// max
	maxClients int32   = 15
	maxCode    int32   = 12345678
	maxCodeLeg int32   = 200
	maxVal     float32 = 10_500_000.00
	maxValLeg  float32 = 865_000.00
	maxFees            = []float32{0.5, 0.8, 0.7, 1.0, 1.2, 1.5}
	// min
	minClients int32   = 200
	minCode    int32   = 23456789
	minCodeLeg int32   = 200
	minVal     float32 = 80_000.00
	minValLeg  float32 = 1_680.00
	minFee             = []float32{1.5, 1.8, 1.9, 2.5, 1.4}
)

// discount
var (
	// fee range
	subFee     float32 = -1.75
	sumFee     float32 = 2.25
	stdDevProp float32 = 0.1
	// total
	discTotalValue float32 = 750_234_567.21
	discAvgFee     float32 = 2.3
)

// concred
var (
	conccredTotalValue           float32 = 750_234_567.21
	conccredTotalEstablishments  int32   = 7_211_563
	conccredActiveEstablishments int32   = 5_001_564
)

// infresta
var (
	infrestaTotalEstablishments int32 = 7_211_563
	infrestaProp                      = []float32{0.5, 0.3, 0.2}
)

// infrterm
var (
	infretermTerminals int32 = 7_316_222
	infretermProp            = []float32{0.3, 0.5, 0.2}
)

// intercam
var (
	intercamTotalValue float32 = 750_234_567.21
	intercamAvgFee     float32 = 1.8
)

// bins
var (
	cardModels = []string{"P", "C"}
	cardProducts = []string{"3", "4", "5", "6", "7", "8", "10", "11", "13", "17", "31", "32", "33", "34", "35", "36", "37", "38"}
)
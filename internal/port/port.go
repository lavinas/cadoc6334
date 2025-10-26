package port

// report domain interface
type Report interface {
	Validate() error
	GetParsedFile(filename string) (map[string]Report, error)
	GetDB(repo Repository) (map[string]Report, error)
	String() string
}

// repository domain interface
type Repository interface {
	FindAll(dest interface{}, conditions ...interface{}) error
}

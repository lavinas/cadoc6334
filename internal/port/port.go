package port

// report domain interface
type Report interface {
	Validate() error
	GetParsedFile(filename string) (map[string]Report, error)
	GetDB(repo Repository) (map[string]Report, error)
	String() string
	Format() string
	GetName() string
}

// repository domain interface
type Repository interface {
	FindAll(dest interface{}, limit int, offset int, conditions ...interface{}) error
}

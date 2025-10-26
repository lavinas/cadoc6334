package port

// report domain interface
type Report interface {
	GetLoaded() (map[string]Report, error)
	GetParsedFile(filename string) (map[string]Report, error)
	String() string
}

// repository domain interface
type Repository interface {
	FindAll(dest interface{}, conditions ...interface{}) error
}

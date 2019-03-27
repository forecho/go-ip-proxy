package storage

type Storage interface {
	Get(string) []byte
	//Delete(string) bool
	Create(string, string) error
	Close() error
	GetAll() []string
	//GetRandomOne() []byte
}

func NewStorage() (Storage, error) {
	return NewLevelDbStorage("db")
}

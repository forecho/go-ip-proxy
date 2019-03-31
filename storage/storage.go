package storage

type Storage interface {
	Get() string
	Delete(string) bool
	Create(string, string) error
	Close()
	GetAll() []string
	//GetRandomOne() []byte
}

func NewStorage() (Storage, error) {
	return NewLevelDbStorage("db")
}

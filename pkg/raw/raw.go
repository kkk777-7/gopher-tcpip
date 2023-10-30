package raw

type Device interface {
	Name() string
	Address() []byte
	Read(data []byte) (int, error)
	Write(data []byte) (int, error)
	Close() error
}

package net

const (
	DEVICE_FLAG_UP = 0x0001
)

type DeviceType string

type Devicer interface {
	Name() string
	Address() string
	Open() error
	Close() error
	Recv(data []byte) (int, error)
	Send(dType DeviceType, data []byte) (int, error)
	Mtu() int
	HeaderSize() int
	AddrSize() int
	Type() DeviceType
}

type Device struct {
	Devicer
	flag int
}

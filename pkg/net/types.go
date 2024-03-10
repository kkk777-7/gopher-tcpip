package net

const (
	DEVICE_FLAG_UP = 0x0001
)

type DeviceType string

const (
	DUMMYDEVICETYPE DeviceType = "dummy"
	LODEVICETYPE    DeviceType = "lo"
)

type ProtocolType string

const (
	IPPROTOOLTYPE ProtocolType = "ip"
)

type Devicer interface {
	Name() string
	Address() string
	Open() error
	Close() error
	Recv(data []byte) (int, error)
	Send(pType ProtocolType, data []byte) (int, error)
	Mtu() int
	HeaderSize() int
	AddrSize() int
	Type() DeviceType
	Priv() interface{}
}

type Device struct {
	Devicer
	flag int
}

type Protocol struct {
	protocolType ProtocolType
	queue        chan *ProtocolEntry
	handler      func(data []byte, len int, dev *Device)
}

type ProtocolEntry struct {
	device *Device
	len    int
	data   []byte
}

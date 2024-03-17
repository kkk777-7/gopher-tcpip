package net

const (
	DEVICE_FLAG_UP = 0x0001
)

type DeviceType uint16

const (
	DUMMYDEVICETYPE DeviceType = 0x0000
	LODEVICETYPE    DeviceType = 0x0001
)

type ProtocolType uint16

const (
	IPPROTOOLTYPE ProtocolType = 0x0800
)

type DeviceCallbackHandler func(device Devicer, protocol ProtocolType, payload []byte) error

type Devicer interface {
	Type() DeviceType
	Name() string
	Address() string
	Close() error
	Read(data []byte) (int, error)
	RxHandler(frame []byte, cb DeviceCallbackHandler) error
	Tx(proto ProtocolType, data []byte) error
	Mtu() int
	HeaderSize() int
}

type Device struct {
	Devicer
}

type ProtocolRxHandler func(dev *Device, data []byte) error

type Protocol struct {
	Type      ProtocolType
	rxQueue   chan *ProtocolEntry
	rxHandler ProtocolRxHandler
}

type ProtocolEntry struct {
	device *Device
	data   []byte
}

package dummy

import (
	"fmt"

	"github.com/kkk777-7/gopher-tcpip/pkg/net"
)

type Devicer interface {
	Name() string
	Address() string
	Open() error
	Close() error
	Recv(data []byte) (int, error)
	Send(dType net.DeviceType, data []byte) (int, error)
	Mtu() int
	HeaderSize() int
	AddrSize() int
	Type() net.DeviceType
}

type Device struct{}

func (d *Device) Name() string {
	return DEVICENAME
}
func (d *Device) Address() string {
	return ""
}
func (d *Device) Open() error  { return nil }
func (d *Device) Close() error { return nil }
func (d *Device) Recv(data []byte) (int, error) {
	return 0, nil
}
func (d *Device) Send(dType net.DeviceType, data []byte) (int, error) {
	fmt.Printf("%s: dev[%s] send: %s, %d\n", DEVICENAME, d.Name(), dType, len(data))
	return 0, nil
}
func (d *Device) Mtu() int {
	return MTU
}
func (d *Device) HeaderSize() int {
	return 0
}
func (d *Device) AddrSize() int {
	return 0
}
func (d *Device) Type() net.DeviceType {
	return net.DUMMYDEVICETYPE
}

func Init() *net.Device {
	fmt.Printf("%s: initialized\n", DEVICENAME)
	return net.RegisterDevice(&Device{})
}

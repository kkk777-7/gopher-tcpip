package dummy

import (
	"fmt"

	"github.com/kkk777-7/gopher-tcpip/pkg/net"
)

func NewDevice() *Device {
	return &Device{
		name: DEVICENAME,
	}
}

func (d *Device) Type() net.DeviceType {
	return net.DUMMYDEVICETYPE
}

func (d *Device) Name() string {
	return d.name
}

func (d *Device) Address() string {
	return "10.0.0.1"
}

func (d *Device) Close() error { return nil }
func (d *Device) Read(data []byte) (int, error) {
	return 0, nil
}

func (d *Device) RxHandler(frame []byte, cb net.DeviceCallbackHandler) error {
	fmt.Printf("dummy_RxHandler: dev=%s, len=%d\n", d.Name(), len(frame))
	cb(d, net.IPPROTOOLTYPE, frame)
	return nil
}

func (d *Device) Tx(proto net.ProtocolType, data []byte) error {
	fmt.Printf("dummy_Tx: dev=%s type=%x, len=%d\n", d.Name(), proto, len(data))
	return nil
}

func (d *Device) Mtu() int {
	return MTU
}

func (d *Device) HeaderSize() int {
	return 0
}

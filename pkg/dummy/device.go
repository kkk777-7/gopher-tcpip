package dummy

import (
	"fmt"

	"github.com/kkk777-7/gopher-tcpip/pkg/net"
	"github.com/kkk777-7/gopher-tcpip/pkg/platform/linux"
)

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
	fmt.Printf("dummy_Send: dev=%s send: type=%s, len=%d\n", d.Name(), dType, len(data))
	linux.RaiseIrq(linux.INTR_DUMMY)
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
func (d *Device) Priv() interface{} {
	return nil
}

func Init() *net.Device {
	fmt.Println("dummy_Init: initialized")
	dev := net.RegisterDevice(&Device{})
	linux.RequestIrq(linux.INTR_DUMMY, Isr, linux.INTR_IRQ_SHARED, DEVICENAME, dev)

	return dev
}

func Isr(irq int, d interface{}) error {
	fmt.Printf("dummy_Isr: irq=%d, dev=%s\n", irq, d.(net.Devicer).Name())
	return nil
}

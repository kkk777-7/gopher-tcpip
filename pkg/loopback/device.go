package loopback

import (
	"fmt"

	"github.com/kkk777-7/gopher-tcpip/pkg/net"
	"github.com/kkk777-7/gopher-tcpip/pkg/platform/linux"
)

func (d *Device) Name() string {
	return d.name
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
	if len(d.queue) >= LO_QUEUE_LIMIT {
		return 0, fmt.Errorf("lo_Send: queue full")
	}
	entry := LoEntry{
		deviceType: dType,
		data:       data,
	}
	d.queue <- entry

	fmt.Printf("lo_Send: (queue pushed) dev=%s send: type=%s, len=%d\n", d.Name(), dType, len(data))
	linux.RaiseIrq(INTR_LO)
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
	return net.LODEVICETYPE
}
func (d *Device) Priv() interface{} {
	return d
}

func Init() *net.Device {
	fmt.Println("lo_Init: initialized")
	lo := &Device{
		irq:   INTR_LO,
		queue: make(chan LoEntry, LO_QUEUE_LIMIT),
	}
	dev, devName := net.RegisterDevice(lo)
	dev.Priv().(*Device).name = devName
	linux.RequestIrq(INTR_LO, Isr, linux.INTR_IRQ_SHARED, dev.Name(), dev)

	return dev
}

func Isr(irq int, d interface{}) error {
	netDev := d.(*net.Device)
	dev := netDev.Priv().(*Device)
	for {
		if len(dev.queue) == 0 {
			break
		}
		entry := <-dev.queue
		fmt.Printf("lo_Isr: (queue popped) irq=%d, dev=%s, type=%s, len=%d\n", irq, dev.Name(), entry.deviceType, len(entry.data))
		netDev.InputHandler(entry.deviceType, entry.data[:], len(entry.data))
	}
	return nil
}

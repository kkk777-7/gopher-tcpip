package net

import (
	"context"
	"fmt"
	"sync"

	"github.com/kkk777-7/gopher-tcpip/pkg/platform/linux"
)

var devices = sync.Map{}
var protocols = sync.Map{}

func RegisterDevice(d Devicer) (*Device, string) {
	length := 0
	devices.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	devName := fmt.Sprintf("dev%d", length)

	dev := &Device{
		Devicer: d,
	}
	devices.LoadOrStore(devName, dev)
	fmt.Printf("net_RegisterDevice: register dev=%s\n", devName)
	return dev, devName
}

func (d *Device) Run(ctx context.Context, wg *sync.WaitGroup) error {
	go linux.RunIrq(ctx, wg)
	if _, exists := devices.Load(d.Name()); !exists {
		return fmt.Errorf("net_Run: link dev=%s is not found", d.Name())
	}
	if err := deviceOpen(d); err != nil {
		return err
	}
	return nil
}

func (d *Device) Shutdown() error {
	if err := deviceClose(d); err != nil {
		return err
	}
	devices.Delete(d.Name())
	return nil
}

func (d *Device) InputHandler(pType ProtocolType, data []byte, len int) error {
	protocols.Range(func(_, value interface{}) bool {
		protocol := value.(*Protocol)
		if protocol.protocolType == pType {
			entry := &ProtocolEntry{
				device: d,
				data:   data,
				len:    len,
			}
			protocol.queue <- entry
			fmt.Printf("net_InputHandler: (queue pushed) dev=%s input: type=%s, len=%d\n", d.Name(), pType, len)
		}
		return true
	})
	// unsupported protocol type ignored
	return nil
}

func (d *Device) Output(pType ProtocolType, data []byte, len int) error {
	if !isUpDevice(d) {
		return fmt.Errorf("net_Output: dev=%s not opened", d.Name())
	}
	if len > d.Mtu() {
		return fmt.Errorf("net_Output: too long: dev=%s mtu=%d len=%d", d.Name(), d.Mtu(), len)
	}
	fmt.Printf("net_Output: dev=%s output: type=%s, len=%d\n", d.Name(), pType, len)
	if _, err := d.Send(pType, data[:len]); err != nil {
		return fmt.Errorf("net_Output: dev=%s send error: %v", d.Name(), err)
	}
	return nil
}

func deviceOpen(d *Device) error {
	if isUpDevice(d) {
		return fmt.Errorf("net_deviceOpen: dev=%s is already up", d.Name())
	}
	if err := d.Open(); err != nil {
		return err
	}
	d.flag |= DEVICE_FLAG_UP
	fmt.Printf("net_deviceOpen: open dev=%s\n", d.Name())
	return nil
}

func deviceClose(d *Device) error {
	if !isUpDevice(d) {
		return fmt.Errorf("net_deviceClose: dev=%s is already down", d.Name())
	}
	if err := d.Close(); err != nil {
		return err
	}
	d.flag &= ^DEVICE_FLAG_UP
	fmt.Printf("net_deviceClose: close dev=%s\n", d.Name())
	return nil
}

func isUpDevice(d *Device) bool {
	return d.flag&DEVICE_FLAG_UP != 0
}

func ResisterProtocol(protocolType ProtocolType, handler func(data []byte, len int, dev *Device)) {
	protocol := &Protocol{
		protocolType: protocolType,
		queue:        make(chan *ProtocolEntry),
		handler:      handler,
	}
	protocols.LoadOrStore(protocol.protocolType, protocol)
	fmt.Printf("net_RegisterProtocol: register protocol=%s\n", protocol.protocolType)
}

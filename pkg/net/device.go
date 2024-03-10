package net

import (
	"fmt"
	"sync"
)

var devices = sync.Map{}

func RegisterDevice(d Devicer) *Device {
	dev := &Device{
		Devicer: d,
	}
	devices.LoadOrStore(dev.Name(), dev)
	return dev
}

func (d *Device) Run() error {
	if _, exists := devices.Load(d.Name()); !exists {
		return fmt.Errorf("link dev[%s] is not found", d.Name())
	}
	if err := deviceOpen(d); err != nil {
		return err
	}
	return nil
}

func (d *Device) Shutdown() error {
	fmt.Printf("close dev[%s]\n", d.Name())
	if err := deviceClose(d); err != nil {
		return err
	}
	devices.Delete(d.Name())
	return nil
}

func (d *Device) InputHandler(dType DeviceType, data []byte, len int) error {
	fmt.Printf("dev[%s] input: %s, %d\n", d.Name(), dType, len)
	return nil
}

func (d *Device) Output(dType DeviceType, data []byte, len int) error {
	if !isUpDevice(d) {
		return fmt.Errorf("dev[%s] not opened", d.Name())
	}
	if len > d.Mtu() {
		return fmt.Errorf("too long: dev[%s] mtu[%d] len[%d]", d.Name(), d.Mtu(), len)
	}
	fmt.Printf("dev[%s] output: %s, %d\n", d.Name(), dType, len)
	if _, err := d.Send(d.Type(), data[:len]); err != nil {
		return fmt.Errorf("dev[%s] send error: %v", d.Name(), err)
	}
	return nil
}

func deviceOpen(d *Device) error {
	if isUpDevice(d) {
		return fmt.Errorf("dev[%s] is already up", d.Name())
	}
	if err := d.Open(); err != nil {
		return err
	}
	d.flag |= DEVICE_FLAG_UP
	fmt.Printf("open dev[%s]\n", d.Name())
	return nil
}

func deviceClose(d *Device) error {
	if !isUpDevice(d) {
		return fmt.Errorf("dev[%s] is already down", d.Name())
	}
	if err := d.Close(); err != nil {
		return err
	}
	d.flag &= ^DEVICE_FLAG_UP
	fmt.Printf("close dev[%s]\n", d.Name())
	return nil
}

func isUpDevice(d *Device) bool {
	return d.flag&DEVICE_FLAG_UP != 0
}

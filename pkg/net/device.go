package net

import (
	"fmt"
	"io"
	"log"
	"sync"
)

type Devicer interface {
	Name() string
	Address() string
	Recv(data []byte) (int, error)
	Send(data []byte) (int, error)
	Close() error
	Handle(data []byte) error
	Mtu() int
	HeaderSize() int
}

type Device struct {
	Devicer
	errors chan error
}

var devices = sync.Map{}

func RegisterDevice(d Devicer) (*Device, error) {
	if _, exists := devices.Load(d); exists {
		return nil, fmt.Errorf("link device '%s' is already registered", d.Name())
	}
	dev := &Device{
		Devicer: d,
		errors:  make(chan error),
	}
	devices.Store(d, dev)
	return dev, nil
}

func (d *Device) Run() error {
	if _, exists := devices.Load(d); exists {
		return fmt.Errorf("link device '%s' is not found", d.Name())
	}
	go func() {
		var buf = make([]byte, d.HeaderSize()+d.Mtu())
		for {
			n, err := d.Recv(buf)
			if n > 0 {
				err := d.Handle(buf)
				if err != nil {
					d.errors <- err
					break
				}
			}
			if err != nil {
				d.errors <- err
				break
			}
		}
		close(d.errors)
	}()
	return nil
}

func (d *Device) Shutdown() error {
	if err := d.Close(); err != nil {
		return err
	}
	close(d.errors)
	if err := <-d.errors; err != nil {
		if err != io.EOF {
			log.Println(err)
		}
	}
	devices.Delete(d.Devicer)
	return nil
}

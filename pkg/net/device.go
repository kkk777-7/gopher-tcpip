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

	go func() {
		var buf = make([]byte, d.HeaderSize()+d.Mtu())
		for {
			n, err := d.Recv(buf)
			if n > 0 {
				err := d.Handle(buf)
				if err != nil {
					dev.errors <- err
					break
				}
			}
			if err != nil {
				dev.errors <- err
				break
			}
		}
		close(dev.errors)
	}()
	devices.Store(d, dev)

	check, _ := devices.Load(d)
	fmt.Printf("check: %v\n", check)
	return dev, nil
}

func (d *Device) Shutdown() {
	d.Devicer.Close()
	if err := <-d.errors; err != nil {
		if err != io.EOF {
			log.Println(err)
		}
	}
	devices.Delete(d.Devicer)
}

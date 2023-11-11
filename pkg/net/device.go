package net

import (
	"context"
	"fmt"
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
	*canceler
}

type canceler struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         *sync.WaitGroup
}

var devices = sync.Map{}

func RegisterDevice(d Devicer) *Device {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	dev := &Device{
		Devicer: d,
		errors:  make(chan error),
		canceler: &canceler{
			ctx:        ctx,
			cancelFunc: cancel,
			wg:         &wg,
		},
	}
	devices.LoadOrStore(dev.Name(), dev)
	return dev
}

func (d *Device) Run() error {
	if _, exists := devices.Load(d.Name()); !exists {
		return fmt.Errorf("link device '%s' is not found", d.Name())
	}

	go func() {
		buf := make([]byte, d.HeaderSize()+d.Mtu())
		for {
			select {
			case <-d.ctx.Done():
				log.Printf("dev[%s] is canceled", d.Name())
				close(d.errors)
				d.wg.Done()
				return
			default:
				n, err := d.Recv(buf)
				if err != nil {
					d.errors <- err
					return
				}
				if n > 0 {
					if err := d.Handle(buf); err != nil {
						d.errors <- err
						return
					}
				}
			}
		}
	}()
	return nil
}

func (d *Device) Shutdown() error {
	d.wg.Add(1)
	d.cancelFunc()
	d.wg.Wait()

	if err := d.Close(); err != nil {
		return err
	}
	err, ok := <-d.errors
	if ok {
		log.Printf("dev[%s] error: %s", d.Name(), err.Error())
		return err
	}
	devices.Delete(d.Devicer)
	return nil
}

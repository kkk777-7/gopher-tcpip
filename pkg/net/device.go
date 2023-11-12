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
	Handle(data []byte, cbFn DeviceCallbackFn) error
	Mtu() int
	HeaderSize() int
}

type DeviceCallbackFn func(device *Device, protocol EthernetType, payload []byte, src, dst HardwareAddress) error

type Device struct {
	Devicer
	protocols *sync.Map
	errors    chan error
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
	pMap := sync.Map{}

	dev := &Device{
		Devicer:   d,
		protocols: &pMap,
		errors:    make(chan error),
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

	// Start protocol handler
	var protocolEntries []*Protocol
	d.protocols.Range(func(k interface{}, v interface{}) bool {
		protocolEntry := v.(*Protocol)
		runProtocolHandler(protocolEntry)
		protocolEntries = append(protocolEntries, protocolEntry)
		return true
	})

	go func() {
		buf := make([]byte, d.HeaderSize()+d.Mtu())
		for {
			select {
			case <-d.ctx.Done():
				log.Printf("dev[%s] is canceled", d.Name())
				close(d.errors)
				for _, pe := range protocolEntries {
					close(pe.Queue)
				}
				d.wg.Done()
				return
			default:
				n, err := d.Recv(buf)
				if err != nil {
					d.errors <- err
					return
				}
				if n > 0 {
					if err := d.Handle(buf[:n], pathThroughProtocol); err != nil {
						d.errors <- err
						return
					}
				}
			}
		}
	}()
	return nil
}

func pathThroughProtocol(device *Device, ethType EthernetType, payload []byte, src, dst HardwareAddress) error {
	var err error
	device.protocols.Range(func(k interface{}, v interface{}) bool {
		var (
			_type      = k.(EthernetType)
			protoEntry = v.(Protocol)
		)
		if ethType == _type {
			dev, exists := devices.Load(device.Name())
			if !exists {
				err = fmt.Errorf("link device '%s' is not found", dev.(*Device).Name())
				return false
			}
			protoEntry.Queue <- &packet{
				dev:  dev.(*Device),
				data: payload,
				src:  src,
				dst:  dst,
			}
			return false
		}
		return true
	})
	return err
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

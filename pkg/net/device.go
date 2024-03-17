package net

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

var devices = sync.Map{}
var protocols = sync.Map{}

func RegisterDevice(ctx context.Context, wg *sync.WaitGroup, device Devicer) (*Device, error) {
	if _, exists := devices.Load(device.Name()); exists {
		return nil, fmt.Errorf("device=%s already registered", device.Name())
	}
	dev := &Device{
		Devicer: device,
	}

	// start rx loop
	go dev.Start(ctx, wg)

	devices.Store(device.Name(), dev)
	fmt.Printf("net_RegisterDevice: register dev=%s\n", dev.Name())
	return dev, nil
}

func (dev *Device) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	var terminate bool
	for !terminate {
		select {
		case <-ctx.Done():
			terminate = true
		default:
			var buf = make([]byte, dev.HeaderSize()+dev.Mtu())

			n, err := dev.Read(buf)
			if n > 0 {
				if err := dev.RxHandler(buf[:n], rxHandler); err != nil {
					fmt.Println(err)
				}
			}
			if err != nil {
				fmt.Println(err)
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func rxHandler(device Devicer, protocolType ProtocolType, payload []byte) error {
	// unsupported protocol type ignored
	protocols.Range(func(key, value interface{}) bool {
		pType := key.(ProtocolType)
		protocol := value.(*Protocol)
		if pType == protocolType {
			dev, _ := devices.Load(device)
			protocol.rxQueue <- &ProtocolEntry{
				device: dev.(*Device),
				data:   payload,
			}
			fmt.Printf("net_rxHandler: (queue pushed) dev=%s input: type=%x, len=%d\n", dev.(*Device).Name(), pType, len(payload))
			return false
		}
		return true
	})
	return nil
}

func (d *Device) Shutdown() error {
	if err := d.Close(); err != nil {
		return err
	}
	fmt.Printf("net_Shutdown: close dev=%s\n", d.Name())
	devices.Delete(d.Name())
	return nil
}

func ResisterProtocol(ctx context.Context, wg *sync.WaitGroup, Type ProtocolType, rxHandler ProtocolRxHandler) {
	if _, ok := protocols.Load(Type); ok {
		log.Printf("net_RegisterProtocol: protocol=%x already registered\n", Type)
		return
	}
	protocol := &Protocol{
		Type:      Type,
		rxQueue:   make(chan *ProtocolEntry),
		rxHandler: rxHandler,
	}
	// start rx loop
	go func() {
		defer wg.Done()

		var terminate bool
		for !terminate {
			select {
			case <-ctx.Done():
				terminate = true
			default:
				for entry := range protocol.rxQueue {
					if err := protocol.rxHandler(entry.device, entry.data); err != nil {
						log.Println(err)
					}
				}
			}
		}
	}()

	protocols.Store(protocol.Type, protocol)
	fmt.Printf("net_RegisterProtocol: register protocol=%x\n", protocol.Type)
}

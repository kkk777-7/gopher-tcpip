package net

import "log"

type Protocol struct {
	Type    EthernetType
	Handler ProtocolHandler
	Queue   chan *packet
}

type ProtocolHandler func(dev *Device, data []byte, src, dst HardwareAddress) error

type packet struct {
	dev  *Device
	data []byte
	src  HardwareAddress
	dst  HardwareAddress
}

func (dev *Device) ResisterProtocol(Type EthernetType, handler ProtocolHandler) {
	protocol := Protocol{
		Type:    Type,
		Handler: handler,
		Queue:   make(chan *packet),
	}
	dev.protocols.LoadOrStore(Type, protocol)
}

func runProtocolHandler(proto *Protocol) {
	go func() {
		for packet := range proto.Queue {
			if err := proto.Handler(packet.dev, packet.data, packet.src, packet.dst); err != nil {
				log.Println(err)
			}
		}
	}()
}

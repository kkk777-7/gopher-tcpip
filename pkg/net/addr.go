package net

type EthernetType uint16

const (
	EthernetIP  EthernetType = 0x0800
	EthernetARP EthernetType = 0x0806
)

type HardwareAddress interface {
	Len() uint8
	ByteAddr() []byte
	StringAddr() string
}

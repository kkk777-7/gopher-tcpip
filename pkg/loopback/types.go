package loopback

import "github.com/kkk777-7/gopher-tcpip/pkg/net"

const (
	DEVICENAME     = "lo"
	MTU            = 65535
	LO_QUEUE_LIMIT = 16
	INTR_LO        = 0x0002
)

type Device struct {
	irq   int
	queue chan LoEntry
}

type LoEntry struct {
	deviceType net.DeviceType
	data       []byte
}

package loopback

import "github.com/kkk777-7/gopher-tcpip/pkg/net"

const (
	DEVICENAME = "lo"
	MTU        = 65535
)

type Device struct {
	name  string
	queue chan []byte
}

type header struct {
	Type net.ProtocolType
}

// type address struct{}

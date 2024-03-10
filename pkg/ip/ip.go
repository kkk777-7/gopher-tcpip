package ip

import (
	"fmt"

	"github.com/kkk777-7/gopher-tcpip/pkg/net"
)

func Init() {
	net.ResisterProtocol(net.IPPROTOOLTYPE, InputHandler)
}

func InputHandler(data []byte, len int, dev *net.Device) {
	fmt.Printf("ip_InputHandler: dev=%s, len=%d\n", dev.Name(), len)
}

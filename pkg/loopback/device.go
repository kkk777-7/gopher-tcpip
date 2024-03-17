package loopback

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/kkk777-7/gopher-tcpip/pkg/net"
)

func NewDevice() *Device {
	return &Device{
		name:  DEVICENAME,
		queue: make(chan []byte),
	}
}

func (d *Device) Type() net.DeviceType {
	return net.LODEVICETYPE
}

func (d *Device) Name() string {
	return d.name
}

func (d *Device) Address() string {
	return "10.0.0.1"
}

func (d *Device) Close() error {
	close(d.queue)
	return nil
}

func (d *Device) Read(buf []byte) (int, error) {
	var err error
	data, ok := <-d.queue
	if !ok {
		err = io.EOF
	}
	return copy(buf, data), err
}

func (d *Device) RxHandler(frame []byte, callback net.DeviceCallbackHandler) error {
	fmt.Printf("loopback_RxHandler: dev=%s, len=%d\n", d.Name(), len(frame))
	hdr := header{}
	buf := bytes.NewBuffer(frame)
	if err := binary.Read(buf, binary.BigEndian, &hdr); err != nil {
		return err
	}
	callback(d, net.IPPROTOOLTYPE, buf.Bytes())
	return nil
}

func (d *Device) Tx(Type net.ProtocolType, data []byte) error {
	fmt.Printf("loopback_Tx: dev=%s, len=%d\n", d.Name(), len(data))
	buf := make([]byte, 2+len(data))
	binary.BigEndian.PutUint16(buf[0:2], uint16(Type))
	copy(buf[2:], data)
	d.queue <- buf
	return nil
}

func (d *Device) Mtu() int {
	return MTU
}

func (d *Device) HeaderSize() int {
	return 0
}

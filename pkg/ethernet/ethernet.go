package ethernet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"reflect"

	"github.com/kkk777-7/gopher-tcpip/pkg/raw"
)

const (
	headerSize     = 14
	maxPayloadSize = 1500
	minPayloadSize = 46
	minFrameSize   = headerSize + minPayloadSize
	maxFrameSize   = headerSize + maxPayloadSize
)

type Device struct {
	raw   raw.Device
	mtu   int
	hSize int
}

func NewDevice(dev raw.Device) (*Device, error) {
	if dev == nil {
		return nil, fmt.Errorf("must input raw device")
	}
	ethDev := &Device{
		raw:   dev,
		mtu:   maxPayloadSize,
		hSize: headerSize,
	}
	return ethDev, nil
}

func (d *Device) Name() string {
	return d.raw.Name()
}
func (d *Device) Address() string {
	return ToStringFromByte(d.raw.Address())
}
func (d *Device) Recv(data []byte) (int, error) {
	return d.raw.Read(data)
}
func (d *Device) Send(data []byte) (int, error) {
	return 0, nil
}
func (d *Device) Close() error {
	return d.raw.Close()
}
func (d *Device) Mtu() int {
	return d.mtu
}
func (d *Device) HeaderSize() int {
	return d.hSize
}

func (d *Device) Handle(data []byte) error {
	log.Println("packet handling start...")
	for {
		ethFrame, err := parse(data)
		if err != nil {
			log.Printf("%s error parse: %v", d.raw.Name(), err)
		}
		if reflect.DeepEqual(ethFrame.Header.Dst, d.raw.Address()) {
			log.Println("ok")
		}
	}
}

func parse(data []byte) (*Frame, error) {
	f := &Frame{}
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &f.Header); err != nil {
		return nil, err
	}
	f.Payload = buf.Bytes()
	return f, nil
}

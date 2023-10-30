package tuntap

import (
	"io"
	"runtime"
)

const device = "/dev/net/tun"

type Tap struct {
	io.ReadWriteCloser
	name    string
	address []byte
}

func NewTap(name string) (*Tap, error) {
	t := &Tap{}
	os := runtime.GOOS
	switch os {
	case "linux":
		n, f, err := openTapInLinux(name)
		if err != nil {
			return nil, err
		}
		addr, err := getAddress(name)
		if err != nil {
			return nil, err
		}
		t.ReadWriteCloser, t.name, t.address = f, n, addr
	}
	return t, nil
}

func (t *Tap) Name() string {
	return t.name
}
func (t *Tap) Address() []byte {
	return t.address
}

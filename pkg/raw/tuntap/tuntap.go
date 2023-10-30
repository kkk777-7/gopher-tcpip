package tuntap

import (
	"os"
	"runtime"
)

const device = "/dev/net/tun"

type Tap struct {
	file    *os.File
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
		t.file, t.name, t.address = f, n, addr
	}
	return t, nil
}

func (t *Tap) Name() string {
	return t.name
}
func (t *Tap) Address() []byte {
	return t.address
}
func (t *Tap) Read(data []byte) (int, error) {
	return t.file.Read(data)
}
func (t *Tap) Write(data []byte) (int, error) {
	return t.file.Write(data)
}
func (t *Tap) Close() error {
	return t.file.Close()
}

package tuntap

import (
	"fmt"
	"net"
	"os"
	"strings"
	"unsafe"

	"github.com/kkk777-7/gopher-tcpip/pkg/utils"
	"golang.org/x/sys/unix"
)

const (
	cIFFTAP  = 0x0002
	cIFFNOPI = 0x1000
)

type ifReq struct {
	Name  [0x10]byte
	Flags uint16
	pad   [0x28 - 0x10 - 2]byte
}

func openTapInLinux(name string) (string, *os.File, error) {
	if len(name) >= unix.IFNAMSIZ {
		return "", nil, fmt.Errorf("name is too long")
	}
	fd, err := unix.Open(device, os.O_RDWR|unix.O_NONBLOCK, 0)
	if err != nil {
		return "", nil, err
	}
	var flags uint16 = cIFFNOPI | cIFFTAP

	if name, err = createInterface(uintptr(fd), name, flags); err != nil {
		return "", nil, err
	}
	if err := utils.ExecCmd("sudo", []string{"ip", "link", "set", "dev", name, "up"}); err != nil {
		return "", nil, err
	}
	return name, os.NewFile(uintptr(fd), "tun"), nil
}

func createInterface(fd uintptr, ifName string, flags uint16) (iFName string, err error) {
	var req ifReq
	req.Flags = flags
	copy(req.Name[:], ifName)

	err = utils.Ioctl(fd, unix.TUNSETIFF, uintptr(unsafe.Pointer(&req)))
	if err != nil {
		return
	}

	iFName = strings.Trim(string(req.Name[:]), "\x00")
	return
}

func getAddress(name string) ([]byte, error) {
	addr, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}
	return addr.HardwareAddr, nil
}

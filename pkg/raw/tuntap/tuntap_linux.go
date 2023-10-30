package tuntap

import (
	"fmt"
	"os"
	"syscall"

	"github.com/kkk777-7/gopher-tcpip/pkg/ioctl"
	"golang.org/x/sys/unix"
)

func openTapInLinux(name string) (string, *os.File, error) {
	if len(name) >= unix.IFNAMSIZ {
		return "", nil, fmt.Errorf("name is too long")
	}
	file, err := os.Open(device)
	if err != nil {
		return "", nil, err
	}
	name, err = ioctl.TUNSETIFF(file.Fd(), name)
	if err != nil {
		return "", nil, err
	}
	// NOTE: tun fd won't work epoll, so set NonBlocking
	// ref: https://github.com/golang/go/issues/38618
	if err := unix.SetNonblock(int(file.Fd()), true); err != nil {
		return "", nil, err
	}
	file = os.NewFile(file.Fd(), device)
	flags, err := ioctl.SIOCGIFFLAGS(name)
	if err != nil {
		file.Close()
		return "", nil, err
	}
	flags |= (syscall.IFF_UP | syscall.IFF_RUNNING)
	if err := ioctl.SIOCSIFFLAGS(name, flags); err != nil {
		file.Close()
		return "", nil, err
	}
	return name, file, nil
}

func getAddress(name string) ([]byte, error) {
	addr, err := ioctl.SIOCGIFHWADDR(name)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

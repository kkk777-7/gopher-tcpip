package ioctl

import (
	"bytes"

	"unsafe"

	"golang.org/x/sys/unix"
)

func SIOCGIFFLAGS(name string) (uint16, error) {
	soc, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return 0, err
	}
	defer unix.Close(soc)
	ifreq := struct {
		name  [unix.IFNAMSIZ]byte
		flags uint16
		_pad  [22]byte
	}{}
	copy(ifreq.name[:unix.IFNAMSIZ-1], name)
	if _, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(soc), unix.SIOCGIFFLAGS, uintptr(unsafe.Pointer(&ifreq))); errno != 0 {
		return 0, errno
	}
	return ifreq.flags, nil
}

func SIOCSIFFLAGS(name string, flags uint16) error {
	soc, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}
	defer unix.Close(soc)
	ifreq := struct {
		name  [unix.IFNAMSIZ]byte
		flags uint16
		_pad  [22]byte
	}{}
	copy(ifreq.name[:unix.IFNAMSIZ-1], name)
	ifreq.flags = flags
	if _, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(soc), unix.SIOCSIFFLAGS, uintptr(unsafe.Pointer(&ifreq))); errno != 0 {
		return errno
	}
	return nil
}

type sockaddr struct {
	family uint16
	addr   [14]byte
}

func SIOCGIFHWADDR(name string) ([]byte, error) {
	soc, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return nil, err
	}
	defer unix.Close(soc)
	ifreq := struct {
		name [unix.IFNAMSIZ]byte
		addr sockaddr
		_pad [8]byte
	}{}
	copy(ifreq.name[:unix.IFNAMSIZ-1], name)
	if _, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(soc), unix.SIOCGIFHWADDR, uintptr(unsafe.Pointer(&ifreq))); errno != 0 {
		return nil, errno
	}
	return ifreq.addr.addr[:], nil
}

func TUNSETIFF(fd uintptr, name string) (string, error) {
	ifreq := struct {
		name  [unix.IFNAMSIZ]byte
		flags uint16
		_pad  [22]byte
	}{}
	copy(ifreq.name[:unix.IFNAMSIZ-1], []byte(name))
	ifreq.flags = unix.IFF_TAP | unix.IFF_NO_PI
	if _, _, errno := unix.Syscall(unix.SYS_IOCTL, fd, unix.TUNSETIFF, uintptr(unsafe.Pointer(&ifreq))); errno != 0 {
		return "", errno
	}
	return string(ifreq.name[:bytes.IndexByte(ifreq.name[:], 0)]), nil
}

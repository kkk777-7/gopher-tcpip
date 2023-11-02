package utils

import (
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

func Ioctl(fd uintptr, request uintptr, argp uintptr) error {
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, fd, uintptr(request), argp)
	if errno != 0 {
		return os.NewSyscallError("ioctl", errno)
	}
	return nil
}

func ExecCmd(command string, args []string) error {
	cmdObj := exec.Command(command, args...)
	cmdObj.Stdout, cmdObj.Stderr = os.Stdout, os.Stderr
	if err := cmdObj.Run(); err != nil {
		return err
	}
	return nil
}

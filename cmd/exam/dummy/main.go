package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kkk777-7/gopher-tcpip/pkg/dummy"
	"github.com/kkk777-7/gopher-tcpip/pkg/net"
)

var sig chan os.Signal

func main() {
	dev := setup()
	defer dev.Shutdown()

	if err := dev.Run(); err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		s := <-sig
		log.Printf("sig: %s\n", s)
		done <- struct{}{}
	}()

	for {
		select {
		case <-done:
			return
		default:
			if err := dev.Output(net.DUMMYDEVICETYPE, []byte("hello"), 5); err != nil {
				log.Println(err)
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func setup() *net.Device {
	// signal handling
	sig = make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	return dummy.Init()
}

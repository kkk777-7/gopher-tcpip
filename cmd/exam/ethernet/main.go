package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kkk777-7/gopher-tcpip/pkg/ethernet"
	"github.com/kkk777-7/gopher-tcpip/pkg/net"
	"github.com/kkk777-7/gopher-tcpip/pkg/raw/tuntap"
)

var devName string
var sig chan os.Signal

func init() {
	flag.StringVar(&devName, "name", "", "device name")
}

func setup() (*net.Device, error) {
	flag.Parse()
	if devName == "" {
		fmt.Println("please device name: ./cmd/exam/tuntap --name xxx")
		os.Exit(1)
	}
	// signal handling
	sig = make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	raw, err := tuntap.NewTap(devName)
	if err != nil {
		return nil, err
	}
	eth, err := ethernet.NewDevice(raw)
	if err != nil {
		return nil, err
	}

	dev, err := net.RegisterDevice(eth)
	if err != nil {
		return nil, err
	}
	return dev, nil
}

func main() {
	log.Print("start")
	dev, err := setup()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[%s] %s\n", dev.Name(), dev.Address())

	s := <-sig
	fmt.Printf("sig: %s\n", s)
	dev.Shutdown()
	log.Print("finish")
}

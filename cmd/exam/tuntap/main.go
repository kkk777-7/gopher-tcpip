package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/kkk777-7/gopher-tcpip/pkg/ethernet"
	"github.com/kkk777-7/gopher-tcpip/pkg/raw/tuntap"
)

var devName string

func init() {
	flag.StringVar(&devName, "name", "", "device name")
}

func main() {
	flag.Parse()
	if devName == "" {
		fmt.Println("please device name: ./cmd/exam/tuntap --name xxx")
		os.Exit(1)
	}

	tap, err := tuntap.NewTap(devName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("name=%s, addr=%s\n", tap.Name(), ethernet.StringAddr(tap.Address()))
	buf := make([]byte, 4096)
	for {
		n, err := tap.Read(buf)
		if err != nil {
			panic(err)
		}
		fmt.Printf("--- [%s] incomming %d bytes data ---\n", tap.Name(), n)
		fmt.Printf("%s", hex.Dump(buf[:n]))
	}
}

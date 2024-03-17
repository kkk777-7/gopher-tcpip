package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/kkk777-7/gopher-tcpip/pkg/dummy"
	"github.com/kkk777-7/gopher-tcpip/pkg/net"
)

func main() {
	dummy := setup()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	var wg sync.WaitGroup
	wg.Add(2)

	dev, err := net.RegisterDevice(ctx, &wg, dummy)
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Shutdown()

	fmt.Printf("[%s] %s\n", dev.Name(), dev.Address())

	go Output(ctx, &wg, dev)
	wg.Wait()
}

func Output(ctx context.Context, wg *sync.WaitGroup, dev *net.Device) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("output canceled")
			return
		default:
			if err := dev.Tx(net.IPPROTOOLTYPE, []byte("hello")); err != nil {
				log.Println(err)
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func setup() *dummy.Device {
	return dummy.NewDevice()
}

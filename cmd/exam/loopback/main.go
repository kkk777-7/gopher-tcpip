package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/kkk777-7/gopher-tcpip/pkg/ip"
	"github.com/kkk777-7/gopher-tcpip/pkg/loopback"
	"github.com/kkk777-7/gopher-tcpip/pkg/net"
)

func main() {
	dev := setup()
	defer dev.Shutdown()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(2)

	if err := dev.Run(ctx, &wg); err != nil {
		log.Fatal(err)
	}
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
			if err := dev.Output(net.IPPROTOOLTYPE, []byte("hello"), 5); err != nil {
				log.Println(err)
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func setup() *net.Device {
	ip.Init()
	return loopback.Init()
}

package linux

import (
	"context"
	"fmt"
	"sync"
)

var irqSig = make(chan int, 1)
var irqs = sync.Map{}

const (
	INTR_IRQ_SHARED = 0x0001
	INTR_DUMMY      = 0x0001
)

type irqEntry struct {
	irq     int
	handler func(irq int, dev interface{}) error
	flag    int
	name    string
	device  interface{}
}

func RequestIrq(irq int, handler func(irq int, dev interface{}) error, flag int, name string, device interface{}) error {
	entry := irqEntry{
		irq,
		handler,
		flag,
		name,
		device,
	}
	val, ok := irqs.Load(entry.irq)
	if ok {
		if (val.(*irqEntry).flag&INTR_IRQ_SHARED == 1) || (entry.flag&INTR_IRQ_SHARED == 1) {
			return fmt.Errorf("irq_RequestIrq: conflicts with already registered irq=%d", entry.irq)
		}
	}
	irqs.Store(entry.irq, &entry)
	fmt.Printf("irq_RequestIrq: request irq=%d, flag=%d, name=%s\n", entry.irq, entry.flag, entry.name)
	return nil
}

func RunIrq(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	var terminate bool

	fmt.Println("irq_RunIrq: start...")

	for !terminate {
		select {
		case <-ctx.Done():
			terminate = true
		case i := <-irqSig:
			irqs.Range(func(key, value interface{}) bool {
				entry := value.(*irqEntry)
				if entry.irq == i {
					entry.handler(i, entry.device)
					return true
				}
				return false
			})
		}
	}
	fmt.Println("\nirq_RunIrq: terminated")
}

func RaiseIrq(irq int) {
	irqSig <- irq
}

package leasepool

import (
	"sync"
	"testing"
)

func TestLeasePool(t *testing.T) {
	p := NewPool(5)
	w := &sync.WaitGroup{}
	for i := 0; i < 500; i++ {
		w.Add(1)
		go testLeasePool(i, p.Get(), w)
	}
	w.Wait()

}

func testLeasePool(i int, l Lease, w *sync.WaitGroup) {
	l.Release()
	w.Done()
}

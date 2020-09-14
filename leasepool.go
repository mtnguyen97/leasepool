package leasepool

import (
	"sync"
)

type Lease interface {
	Release()
}

type lease struct {
	r    func()
	once sync.Once
}

func (l *lease) Release() {
	l.once.Do(l.r)
}

type LeasePool interface {
	Get() Lease
}

type leasepool struct {
	remainings int
	cond       *sync.Cond
}

func NewPool(max int) LeasePool {
	return &leasepool{max, sync.NewCond(&sync.Mutex{})}
}

func (p *leasepool) Get() Lease {
	p.cond.L.Lock()
	if p.remainings < 1 {
		p.cond.Wait()
	}
	p.remainings--
	p.cond.L.Unlock()
	return &lease{r: p.release}
}

func (p *leasepool) release() {
	p.cond.L.Lock()
	p.remainings++
	if p.remainings == 1 {
		p.cond.Signal()
	}
	p.cond.L.Unlock()
}

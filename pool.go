package main

import "C"

import (
	"fmt"
	"sync"
)

// ObjectPool keeps track of golang objects and the destructor that is association with
// the C counterpart that was passed to caller.
type ObjectPool struct {
	objs    map[int64]poolEntry
	counter int64
	lock    sync.RWMutex
}

type poolEntry struct {
	GoObject    interface{}
	Desctructor func()
}

func newObjectPool() *ObjectPool {
	return &ObjectPool{
		objs: make(map[int64]poolEntry),
	}
}

func (p *ObjectPool) Get(id int64) any {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.objs[id]
}

func (p *ObjectPool) Add(x interface{}, destructor func()) int64 {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.counter++
	id := p.counter
	p.objs[id] = poolEntry{GoObject: x, Desctructor: destructor}
	return id
}

func (p *ObjectPool) Remove(poolId int64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if e, found := p.objs[poolId]; found {
		e.Desctructor()
		delete(p.objs, poolId)
	} else {
		// start out being fail-fast as this could represent sloppy
		// accounting that should be fixed
		panic(fmt.Sprintf("attempting to free %d that was not found or already freed", poolId))
	}
}

var pool = newObjectPool()

//export fc_free_pool_item
func fc_free_pool_item(poolId C.long) {
	pool.Remove(int64(poolId))
}

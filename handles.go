package lang

import (
	"context"
	"fmt"
	"sync"

	"github.com/freeconf/lang/pb"
)

var Handles = newHandlePool()

// ObjectPool keeps track of golang objects and the destructor that is association with
// the C counterpart that was passed to caller.
type HandlePool struct {
	objects map[uint64]any
	handles map[any]uint64
	counter uint64
	lock    sync.RWMutex
}

type HandleService struct {
	pb.UnimplementedHandlesServer
}

func (s *HandleService) Release(ctx context.Context, in *pb.ReleaseRequest) (*pb.Void, error) {
	Handles.Release(in.GHnd)
	return &pb.Void{}, nil
}

func newHandlePool() *HandlePool {
	return &HandlePool{
		objects: make(map[uint64]any),
		handles: make(map[any]uint64),
		// Go get's even handles, X get's odd so they never collide
		counter: 2,
	}
}

func (p *HandlePool) Reserve() uint64 {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.counter = p.counter + 2
	id := p.counter
	return id
}

func (p *HandlePool) CompleteReservation(reservation uint64, x any) {
	p.lock.Lock()
	defer p.lock.Unlock()
	_, found := p.handles[reservation]
	if found {
		panic(fmt.Sprintf("reservation %d already fullfilled, cannot store %t", reservation, x))
	}
	p.objects[reservation] = x
}

func (p *HandlePool) Require(handle uint64) any {
	p.lock.Lock()
	defer p.lock.Unlock()
	x, found := p.objects[handle]
	if !found {
		panic(fmt.Sprintf("attempting to reference handle %d that was not found", handle))
	}
	return x
}

func (p *HandlePool) Get(handle uint64) any {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.objects[handle]
}

func (p *HandlePool) Hnd(obj any) uint64 {
	p.lock.Lock()
	defer p.lock.Unlock()
	hnd, found := p.handles[obj]
	if !found {
		// start out being fail-fast as this could represent sloppy
		// accounting that should be fixed
		panic(fmt.Sprintf("attempting to reference obj %v that was not found", obj))
	}
	return hnd
}

func (p *HandlePool) Put(x any) uint64 {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.counter++
	hnd := p.counter
	p.objects[hnd] = x
	p.handles[x] = hnd
	return hnd
}

func (p *HandlePool) Record(x any, hnd uint64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.objects[hnd] = x
	p.handles[x] = hnd
}

func (p *HandlePool) Release(handle uint64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if obj, found := p.objects[handle]; found {
		delete(p.objects, handle)
		delete(p.handles, obj)
	} else {
		// start out being fail-fast as this could represent sloppy
		// accounting that should be fixed
		panic(fmt.Sprintf("attempting to release handle %d that was not found or already released", handle))
	}
}

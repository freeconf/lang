package lang

import (
	"context"
	"fmt"
	"sync"

	"github.com/freeconf/lang/pb"
)

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
	d *Driver
}

func (s *HandleService) Release(ctx context.Context, in *pb.ReleaseRequest) (*pb.ReleaseResponse, error) {
	s.d.handles.Release(in.Hnd)
	return &pb.ReleaseResponse{}, nil
}

func newHandlePool() *HandlePool {
	return &HandlePool{
		objects: make(map[uint64]any),
		handles: make(map[any]uint64),
		counter: 100,
	}
}

func (p *HandlePool) Reserve() uint64 {
	p.lock.Lock()
	defer p.lock.Unlock()
	id := p.nextHnd()
	return id
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
		hnd = p.nextHnd()
		p.handles[obj] = hnd
		p.objects[hnd] = obj
	}
	return hnd
}

func (p *HandlePool) Put(x any) uint64 {
	p.lock.Lock()
	defer p.lock.Unlock()
	hnd := p.nextHnd()
	p.objects[hnd] = x
	p.handles[x] = hnd
	return hnd
}

func (p *HandlePool) nextHnd() uint64 {
	next := p.counter
	p.counter = p.counter + 1
	return next
}

func (p *HandlePool) Record(x any, hnd uint64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.objects[hnd] = x
	p.handles[x] = hnd
}

func (p *HandlePool) Release(handle uint64) {
	fmt.Printf("releasing %d\n", handle)
	p.lock.Lock()
	defer p.lock.Unlock()
	if obj, found := p.objects[handle]; found {
		delete(p.objects, handle)
		delete(p.handles, obj)
		// } else {
		// 	// start out being fail-fast as this could represent sloppy
		// 	// accounting that should be fixed
		// 	panic(fmt.Sprintf("attempting to release handle %d that was not found or already released", handle))
	}
}

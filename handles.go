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
	handles map[uint64]any
	counter uint64
	lock    sync.RWMutex
}

type HandleService struct {
	pb.UnimplementedHandlesServer
}

func (s *HandleService) Release(ctx context.Context, in *pb.Handle) (*pb.Void, error) {
	Handles.Release(in.Handle)
	return &pb.Void{}, nil
}

func newHandlePool() *HandlePool {
	return &HandlePool{
		handles: make(map[uint64]any),
	}
}

func (p *HandlePool) Get(handle uint64) any {
	p.lock.Lock()
	defer p.lock.Unlock()
	x, found := p.handles[handle]
	if !found {
		// start out being fail-fast as this could represent sloppy
		// accounting that should be fixed
		panic(fmt.Sprintf("attempting to reference handle %d that was not found", handle))
	}
	return x
}

func (p *HandlePool) Put(x any) uint64 {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.counter++
	id := p.counter
	p.handles[id] = x
	return id
}

func (p *HandlePool) Release(handle uint64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if _, found := p.handles[handle]; found {
		delete(p.handles, handle)
	} else {
		// start out being fail-fast as this could represent sloppy
		// accounting that should be fixed
		panic(fmt.Sprintf("attempting to release handle %d that was not found or already released", handle))
	}
}

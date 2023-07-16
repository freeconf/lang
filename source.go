package lang

import (
	"context"
	"errors"

	"github.com/freeconf/lang/pb"
	"github.com/freeconf/restconf"
	"github.com/freeconf/yang"
	"github.com/freeconf/yang/source"
)

func (s *ParserService) Source(ctx context.Context, in *pb.SourceRequest) (*pb.SourceResponse, error) {
	var opener source.Opener
	if in.YangInternalYpath {
		opener = yang.InternalYPath
	} else if in.RestconfInternalYpath {
		opener = restconf.InternalYPath
	} else if in.Path != "" {
		opener = source.Path(in.Path)
	} else if len(in.Any) > 0 {
		openers := make([]source.Opener, len(in.Any))
		for i := 0; i < len(openers); i++ {
			openers[i] = resolveOpener(s.d.handles, in.Any[i])
		}
		opener = source.Any(openers...)
	} else {
		return nil, errors.New("no valid source criteria given")
	}
	return &pb.SourceResponse{SourceHnd: putOpenener(s.d.handles, opener)}, nil
}

func putOpenener(handles *HandlePool, opener source.Opener) uint64 {
	// In golang, apparently you cannot put a func pointer type into a map, so we'll put pointer
	// as long as we expected that when resolving the handle
	return handles.Put(&opener)
}

func resolveOpener(handles *HandlePool, sourceHnd uint64) source.Opener {
	var ypath source.Opener
	if src := handles.Get(sourceHnd); src != nil {
		ptr := src.(*source.Opener)
		ypath = *ptr
	}
	return ypath
}

package lang

import (
	"context"
	"errors"
	"io"

	"github.com/freeconf/lang/pb"
	"github.com/freeconf/yang/meta"
	"github.com/freeconf/yang/parser"
)

type ParserService struct {
	pb.UnimplementedParserServer
	d *Driver
}

func (s *ParserService) LoadModule(ctx context.Context, in *pb.LoadModuleRequest) (*pb.LoadModuleResponse, error) {
	// In golang, apparently you cannot put a func pointer type into a map, so we put the pointer
	// in and so we need to expected that here when resolving the handle
	ypath := resolveOpener(s.d.handles, in.SourceHnd)
	var m *meta.Module
	var err error
	if in.GetName() != "" {
		m, err = parser.LoadModule(ypath, in.GetName())
		if err != nil {
			return nil, err
		}
	} else if in.GetStreamHnd() != 0 {
		rdr := s.d.handles.Require(in.GetStreamHnd()).(io.Reader)
		data, err := io.ReadAll(rdr)
		if err != nil {
			return nil, err
		}
		m, err = parser.LoadModuleFromString(ypath, string(data))
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("must supply either stream or module name")
	}
	return &pb.LoadModuleResponse{
		ModuleHnd: s.d.handles.Put(m),
		Module:    NewMetaEncoder().Encode(m),
	}, nil
}

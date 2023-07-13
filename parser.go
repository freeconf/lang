package lang

import (
	"context"
	"errors"
	"io"

	"github.com/freeconf/lang/pb"
	"github.com/freeconf/yang/meta"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

type ParserService struct {
	pb.UnimplementedParserServer
	d *Driver
}

func (s *ParserService) LoadModule(ctx context.Context, in *pb.LoadModuleRequest) (*pb.LoadModuleResponse, error) {
	var ypath source.Opener
	if in.Path != "" {
		ypath = source.Path(in.Path)
	}
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
		Module:    new(MetaEncoder).Encode(m),
	}, nil
}

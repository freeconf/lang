package lang

import (
	"context"

	"github.com/freeconf/lang/pb"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

type ParserService struct {
	pb.UnimplementedParserServer
}

func (s *ParserService) LoadModule(ctx context.Context, in *pb.LoadModuleRequest) (*pb.LoadModuleResponse, error) {
	ypath := source.Dir(in.Dir)
	m, err := parser.LoadModule(ypath, in.Name)
	if err != nil {
		return nil, err
	}
	return &pb.LoadModuleResponse{
		ModuleHnd: Handles.Put(m),
		Module:    new(MetaEncoder).Encode(m),
	}, nil
}

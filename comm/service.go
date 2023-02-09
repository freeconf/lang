package comm

import (
	"context"

	"github.com/freeconf/lang/comm/pb"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

type Service struct {
	pb.UnimplementedParserServer
}

func (s *Service) LoadModule(ctx context.Context, in *pb.LoadModuleRequest) (*pb.Module, error) {
	ypath := source.Dir(in.Dir)
	m, err := parser.LoadModule(ypath, in.Name)
	if err != nil {
		return nil, err
	}
	return new(MetaEncoder).Encode(m), nil
}

package comm

import (
	"context"

	"github.com/freeconf/lang/comm/pb"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

type Service struct {
	pb.UnimplementedParserServer
	pb.UnimplementedDriverServer
}

func (s *Service) LoadModule(ctx context.Context, in *pb.LoadModuleRequest) (*pb.LoadModuleResponse, error) {
	ypath := source.Dir(in.Dir)
	m, err := parser.LoadModule(ypath, in.Name)
	if err != nil {
		return nil, err
	}
	return &pb.LoadModuleResponse{
		Handle: Handles.Put(m),
		Module: new(MetaEncoder).Encode(m),
	}, nil
}

func (s *Service) Release(ctx context.Context, in *pb.ReleaseRequest) (*pb.Void, error) {
	Handles.Release(in.Handle)
	return &pb.Void{}, nil
}

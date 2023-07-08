package lang

import (
	"context"
	"io"

	"github.com/freeconf/lang/pb"
	"github.com/freeconf/yang/nodeutil"
)

type NodeUtilService struct {
	pb.UnimplementedNodeUtilServer
	d *Driver
}

func (s *NodeUtilService) JSONRdr(ctx context.Context, req *pb.JSONRdrRequest) (*pb.JSONRdrResponse, error) {
	rdr := s.d.handles.Require(req.StreamHnd).(io.Reader)
	jrdr := nodeutil.ReadJSONIO(rdr)
	return &pb.JSONRdrResponse{NodeHnd: s.d.handles.Put(jrdr)}, nil
}

func (s *NodeUtilService) JSONWtr(ctx context.Context, req *pb.JSONWtrRequest) (*pb.JSONWtrResponse, error) {
	wtr := s.d.handles.Require(req.StreamHnd).(io.Writer)
	jwtr := nodeutil.NewJSONWtr(wtr)
	return &pb.JSONWtrResponse{NodeHnd: s.d.handles.Put(jwtr.Node())}, nil
}

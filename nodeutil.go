package lang

import (
	"context"
	"os"

	"github.com/freeconf/lang/pb"
	"github.com/freeconf/yang/nodeutil"
)

type NodeUtilService struct {
	pb.UnimplementedNodeUtilServer
	d *Driver
}

func (s *NodeUtilService) JSONRdr(ctx context.Context, req *pb.JSONRdrRequest) (*pb.JSONRdrResponse, error) {
	rdr, err := openFileHandle(ctx, s.d, req.File)
	if err != nil {
		return nil, err
	}
	defer rdr.Close()
	jrdr := nodeutil.ReadJSONIO(rdr)
	return &pb.JSONRdrResponse{NodeHnd: s.d.handles.Put(jrdr)}, nil
}

func (s *NodeUtilService) JSONWtr(ctx context.Context, req *pb.JSONWtrRequest) (*pb.JSONWtrResponse, error) {
	f, err := os.Create(req.Fname)
	if err != nil {
		return nil, err
	}
	wtr := nodeutil.NewJSONWtr(f)
	return &pb.JSONWtrResponse{NodeHnd: s.d.handles.Put(wtr.Node())}, nil
}

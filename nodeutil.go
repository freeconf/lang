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
	f, err := os.Open(req.Fname)
	if err != nil {
		return nil, err
	}
	rdr := nodeutil.ReadJSONIO(f)
	return &pb.JSONRdrResponse{NodeHnd: s.d.handles.Put(rdr)}, nil
}

package lang

import (
	"context"
	"os"

	"github.com/freeconf/lang/pb"
	"github.com/freeconf/yang/nodeutil"
)

type NodeUtilService struct {
	pb.UnimplementedNodeUtilServer
}

func (*NodeUtilService) JSONRdr(ctx context.Context, req *pb.JSONRdrRequest) (*pb.JSONRdrResponse, error) {
	f, err := os.Open(req.Fname)
	if err != nil {
		return nil, err
	}
	rdr := nodeutil.ReadJSONIO(f)
	return &pb.JSONRdrResponse{GNodeHnd: Handles.Put(rdr)}, nil
}

package lang

import (
	"context"

	"github.com/freeconf/lang/pb"
	"github.com/freeconf/restconf"
	"github.com/freeconf/restconf/device"
)

type RestconfService struct {
	pb.UnimplementedRestconfServer
	d *Driver
}

func (s *RestconfService) NewServer(ctx context.Context, req *pb.NewServerRequest) (*pb.NewServerResponse, error) {
	d := s.d.handles.Require(req.DeviceHnd).(*device.Local)
	server := restconf.NewServer(d)
	resp := pb.NewServerResponse{
		ServerHnd: s.d.handles.Put(server),
	}
	return &resp, nil
}

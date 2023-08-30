package lang

import (
	"context"

	"github.com/freeconf/gnmi"
	"github.com/freeconf/lang/pb"
	"github.com/freeconf/restconf"
	"github.com/freeconf/restconf/device"
)

type ProtoService struct {
	pb.UnimplementedProtoServer
	d *Driver
}

func (s *ProtoService) RestconfServer(ctx context.Context, req *pb.RestconfServerRequest) (*pb.RestconfServerResponse, error) {
	d := s.d.handles.Require(req.DeviceHnd).(*device.Local)
	server := restconf.NewServer(d)
	resp := pb.RestconfServerResponse{
		ServerHnd: s.d.handles.Put(server),
	}
	return &resp, nil
}

func (s *ProtoService) GnmiServer(ctx context.Context, req *pb.GnmiServerRequest) (*pb.GnmiServerResponse, error) {
	d := s.d.handles.Require(req.DeviceHnd).(*device.Local)
	server := gnmi.NewServer(d)
	resp := pb.GnmiServerResponse{
		ServerHnd: s.d.handles.Put(server),
	}
	return &resp, nil
}

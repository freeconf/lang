package lang

import (
	"context"
	"io"

	"github.com/freeconf/lang/pb"
	"github.com/freeconf/restconf/client"
	"github.com/freeconf/restconf/device"
	"github.com/freeconf/yang/node"
)

type DeviceService struct {
	pb.UnimplementedDeviceServer
	d *Driver
}

func (s *DeviceService) NewDevice(ctx context.Context, req *pb.NewDeviceRequest) (*pb.NewDeviceResponse, error) {
	ypath := resolveOpener(s.d.handles, req.YangPathSourceHnd)
	d := device.New(ypath)
	resp := pb.NewDeviceResponse{
		DeviceHnd: s.d.handles.Put(d),
	}
	return &resp, nil
}

func (s *DeviceService) DeviceAddBrowser(ctx context.Context, req *pb.DeviceAddBrowserRequest) (*pb.DeviceAddBrowserResponse, error) {
	d := s.d.handles.Require(req.DeviceHnd).(*device.Local)
	b := s.d.handles.Require(req.BrowserHnd).(*node.Browser)
	d.AddBrowser(b)
	return &pb.DeviceAddBrowserResponse{}, nil
}

func (s *DeviceService) DeviceGetBrowser(ctx context.Context, req *pb.DeviceGetBrowserRequest) (*pb.DeviceGetBrowserResponse, error) {
	d := s.d.handles.Require(req.DeviceHnd).(device.Device)
	b, err := d.Browser(req.ModuleIdent)
	if err != nil {
		return nil, err
	}
	resp := pb.DeviceGetBrowserResponse{
		BrowserHnd: s.d.handles.Hnd(b),
	}
	return &resp, nil
}

func (s *DeviceService) ApplyStartupConfig(ctx context.Context, req *pb.ApplyStartupConfigRequest) (*pb.ApplyStartupConfigResponse, error) {
	d := s.d.handles.Require(req.DeviceHnd).(*device.Local)
	rdr := s.d.handles.Require(req.StreamHnd).(io.Reader)
	err := d.ApplyStartupConfig(rdr)
	resp := pb.ApplyStartupConfigResponse{}
	return &resp, err
}

func (s *DeviceService) Client(ctx context.Context, req *pb.ClientRequest) (*pb.ClientResponse, error) {
	ypath := resolveOpener(s.d.handles, req.YpathHnd)
	proto := client.ProtocolHandler(ypath)
	dev, err := proto(req.Address)
	if err != nil {
		return nil, err
	}
	resp := pb.ClientResponse{
		DeviceHnd: s.d.handles.Hnd(dev),
	}
	return &resp, nil
}

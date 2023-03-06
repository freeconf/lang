package lang

import (
	"context"

	"github.com/freeconf/lang/pb"
	"github.com/freeconf/restconf/device"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/source"
)

type DeviceService struct {
	pb.UnimplementedDeviceServer
	d *Driver
}

func (s *DeviceService) NewDevice(ctx context.Context, req *pb.NewDeviceRequest) (*pb.NewDeviceResponse, error) {
	ypath := source.Path(req.YangPath)
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
	err := d.ApplyStartupConfigFile(req.ConfigFile)
	resp := pb.ApplyStartupConfigResponse{}
	return &resp, err
}

package lang

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/freeconf/lang/pb"
	"github.com/freeconf/yang/fc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Driver struct {
	listener net.Listener
	gserver  *grpc.Server
	pb.UnimplementedNodeServer
	xnodes      pb.XNodeClient
	handles     *HandlePool
	xclientAddr string
}

func NewDriver(gServerAddr string, xClientAddr string) (*Driver, error) {
	d := &Driver{
		handles:     newHandlePool(),
		xclientAddr: xClientAddr,
	}
	if err := d.createXClient(xClientAddr); err != nil {
		return nil, err
	}
	if err := d.createGServer(gServerAddr); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Driver) createGServer(addr string) error {
	var err error
	if _, ferr := os.Stat(addr); ferr == nil {
		if err = os.Remove(addr); err != nil {
			return fmt.Errorf("could not remove old socket %s file. %w", addr, err)
		}
	}
	if d.listener, err = net.Listen("unix", addr); err != nil {
		return fmt.Errorf("listen error. %s. %w", addr, err)
	}
	fc.Debug.Printf("started server on %s", addr)
	d.gserver = grpc.NewServer()
	pb.RegisterParserServer(d.gserver, &ParserService{d: d})
	pb.RegisterHandlesServer(d.gserver, &HandleService{d: d})
	pb.RegisterNodeServer(d.gserver, &NodeService{d: d})
	pb.RegisterNodeUtilServer(d.gserver, &NodeUtilService{d: d})
	pb.RegisterDeviceServer(d.gserver, &DeviceService{d: d})
	pb.RegisterRestconfServer(d.gserver, &RestconfService{d: d})
	return nil
}

func (s *Driver) createXClient(addr string) error {
	credentials := insecure.NewCredentials()
	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
		var d net.Dialer
		return d.DialContext(ctx, "unix", addr)
	}
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials),
		grpc.WithBlock(),
		grpc.WithContextDialer(dialer),
	}
	channel, err := grpc.Dial(addr, options...)
	if err != nil {
		return fmt.Errorf("failed to start client to x server on %s. %w", addr, err)
	}
	fc.Debug.Printf("connected to %s", addr)
	s.xnodes = pb.NewXNodeClient(channel)
	return nil
}

// Serve is a blocking call that starts the GRPC server
func (s *Driver) Serve() error {
	defer s.listener.Close()
	if err := s.gserver.Serve(s.listener); err != nil {
		return fmt.Errorf("grpc server error. %w", err)
	}
	return nil
}

func (s *Driver) Stop() {
	s.gserver.Stop()
}

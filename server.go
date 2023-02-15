package lang

import (
	"fmt"
	"net"
	"os"

	"github.com/freeconf/lang/pb"
	"google.golang.org/grpc"
)

type Service struct {
	listener net.Listener
	gserver  *grpc.Server
	pb.UnimplementedNodeServer
	Client pb.XNodeClient
}

func NewService(addr string, c pb.XNodeClient) (*Service, error) {
	if _, ferr := os.Stat(addr); ferr == nil {
		if err := os.Remove(addr); err != nil {
			return nil, fmt.Errorf("could not remove old socket file. %w", err)
		}
	}
	l, err := net.Listen("unix", addr)
	if err != nil {
		return nil, fmt.Errorf("listen error. %w", err)
	}
	impl := &Service{
		listener: l,
		gserver:  grpc.NewServer(),
	}
	pb.RegisterParserServer(impl.gserver, &ParserService{})
	pb.RegisterHandlesServer(impl.gserver, &HandleService{})
	pb.RegisterNodeServer(impl.gserver, &NodeService{Client: c})
	return impl, nil
}

func (s *Service) Serve() error {
	defer s.listener.Close()
	if err := s.gserver.Serve(s.listener); err != nil {
		return fmt.Errorf("grpc server error. %w", err)
	}
	return nil
}

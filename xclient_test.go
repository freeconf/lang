package lang

import (
	"context"
	"net"
	"os"
	"testing"
	"time"

	"github.com/freeconf/lang/pb"
	"github.com/freeconf/yang/fc"
	"google.golang.org/grpc"
)

func createGrpcServer(t *testing.T, addr string) (*grpc.Server, net.Listener) {
	if _, ferr := os.Stat(addr); ferr == nil {
		if err := os.Remove(addr); err != nil {
			t.Fatalf("could not remove old socket file. %s", err)
		}
	}
	l, err := net.Listen("unix", addr)
	fc.AssertEqual(t, nil, err)
	<-time.After(100 * time.Millisecond)
	return grpc.NewServer(), l
}

func TestXClient(t *testing.T) {
	addr := "/tmp/foo"
	s, l := createGrpcServer(t, addr)
	pb.RegisterXNodeServer(s, &dummyXNode{})
	defer l.Close()
	go s.Serve(l)
	<-time.After(10 * time.Millisecond)
	c := createGrpcClient(t, addr)
	defer c.Close()
	xc := pb.NewXNodeClient(c)
	ctx := context.Background()
	req := &pb.ChildRequest{}
	resp, err := xc.Child(ctx, req)
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, uint64(1000), resp.Handle)
}

type dummyXNode struct {
	pb.UnimplementedXNodeServer
}

func (*dummyXNode) Child(context.Context, *pb.ChildRequest) (*pb.ChildResponse, error) {
	return &pb.ChildResponse{Handle: 1000}, nil
}

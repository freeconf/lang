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
	"google.golang.org/grpc/credentials/insecure"
)

func TestDriverServer(t *testing.T) {
	addr := "/tmp/foo"
	s, err := NewDriver(addr, "")
	go func() {
		if err = s.Serve(); err != nil {
			panic(err)
		}
	}()
	fc.AssertEqual(t, nil, err)

	c := createGrpcClient(t, addr)
	defer c.Close()

	client := pb.NewParserClient(c)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.LoadModule(ctx, &pb.LoadModuleRequest{Dir: "./test/testdata/yang", Name: "basic"})
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, "basic", resp.Module.Ident)
}

func TestDriverClient(t *testing.T) {
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
	req := &pb.XChildRequest{}
	resp, err := xc.XChild(ctx, req)
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, uint64(1000), resp.NodeHnd)
}

type dummyXNode struct {
	pb.UnimplementedXNodeServer
}

func (*dummyXNode) XChild(context.Context, *pb.XChildRequest) (*pb.XChildResponse, error) {
	return &pb.XChildResponse{NodeHnd: 1000}, nil
}

func createGrpcClient(t *testing.T, addr string) *grpc.ClientConn {
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
	c, err := grpc.Dial(addr, options...)
	fc.AssertEqual(t, nil, err)
	return c
}

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

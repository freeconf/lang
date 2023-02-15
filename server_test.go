package lang

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/freeconf/lang/pb"
	"github.com/freeconf/yang/fc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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

func TestService(t *testing.T) {
	addr := "/tmp/foo"
	s, err := NewService(addr, nil)
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
	resp, err := client.LoadModule(ctx, &pb.LoadModuleRequest{Dir: "./test/yang", Name: "testme-1"})
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, "testme-1", resp.Module.Ident)
}

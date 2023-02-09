package comm

import (
	"context"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/freeconf/lang/comm/pb"
	"github.com/freeconf/yang/fc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func startServer(addr string) {
	if err := os.RemoveAll(addr); err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("unix", addr)
	if err != nil {
		log.Fatal("listen error:", err)
	}

	s := grpc.NewServer()
	pb.RegisterParserServer(s, &Service{})
	go func() {
		defer l.Close()
		if err := s.Serve(l); err != nil {
			log.Fatalf("grpc server error. %s", err)
		}
	}()
}

func TestSerice(t *testing.T) {
	addr := "/tmp/foo"
	startServer(addr)

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
	defer c.Close()

	client := pb.NewParserClient(c)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.LoadModule(ctx, &pb.LoadModuleRequest{Dir: "../test/yang", Name: "testme-1"})
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, "testme-1", resp.Ident)
	select {}
}

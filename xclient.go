package lang

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/freeconf/lang/pb"
)

func CreateXClient(addr string) (pb.XNodeClient, error) {
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
	if err != nil {
		return nil, fmt.Errorf("failed to start client to x server. %w", err)
	}
	return pb.NewXNodeClient(c), nil
}

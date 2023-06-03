package test

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/freeconf/lang"
	"github.com/freeconf/lang/pb"
	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/node"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type x interface {
	connect(gAddr string, xAddr string) error
	stop() error
}

type Harness struct {
	name   string
	driver *lang.Driver
	client pb.TestHarnessClient
	access *lang.TestHarnessAccess
	x      x
	xAddr  string
	gAddr  string
}

func NewHarness(name string, x x) *Harness {
	cwd, _ := os.Getwd()
	return &Harness{
		name:  name,
		x:     x,
		xAddr: cwd + "/fc-test-x.sock",
		gAddr: cwd + "/fc-test.sock",
	}
}

func (h *Harness) Close() error {
	return h.x.stop()
}

func (h *Harness) Name() string {
	return h.name
}

func (h *Harness) Connect() error {
	fc.Debug.Println("starting x server")
	if err := h.x.connect(h.gAddr, h.xAddr); err != nil {
		return err
	}

	fc.Debug.Print("waiting for x server to start")
	for i := 0; i <= 20; i++ {
		if _, ferr := os.Stat(h.xAddr); ferr == nil {
			fc.Debug.Println("")
			break
		}
		if i == 20 {
			return fmt.Errorf("timed out waiting for %s", h.xAddr)
		}
		fc.Debug.Print(".")
		time.Sleep(500 * time.Millisecond)
	}

	fc.Debug.Println("connecting test x client")
	if err := h.ConnectClient(); err != nil {
		return err
	}
	var err error
	fc.Debug.Println("starting g server")
	if h.driver, err = lang.NewDriver(h.gAddr, h.xAddr); err != nil {
		return err
	}
	time.Sleep(time.Second)
	h.access = lang.NewTestHarnessClient(h.driver)
	go h.driver.Serve()
	return nil
}

func (h *Harness) ConnectClient() error {
	credentials := insecure.NewCredentials()
	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
		var d net.Dialer
		return d.DialContext(ctx, "unix", h.xAddr)
	}
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials),
		grpc.WithBlock(),
		grpc.WithContextDialer(dialer),
	}
	channel, err := grpc.Dial(h.xAddr, options...)
	if err != nil {
		return fmt.Errorf("failed to start client to x server on %s. %w", h.xAddr, err)
	}
	fc.Debug.Printf("connected to %s", h.xAddr)
	h.client = pb.NewTestHarnessClient(channel)
	return nil
}

func (h *Harness) Stop() error {
	defer h.driver.Stop()
	return h.x.stop()
}

func (h *Harness) createTestCase(tc pb.TestCase, traceFile string) (node.Node, error) {
	req := pb.CreateTestCaseRequest{TestCase: tc, TraceFile: traceFile}
	resp, err := h.client.CreateTestCase(context.Background(), &req)
	if err != nil {
		return nil, err
	}
	return h.access.CreateNode(resp.NodeHnd), nil
}

func (h *Harness) finalizeTestCase() error {
	req := pb.FinalizeTestCaseRequest{}
	_, err := h.client.FinalizeTestCase(context.Background(), &req)
	return err
}

func (h *Harness) parseModule(dir string, moduleIdent string, dumpFile string) error {
	req := pb.ParseModuleRequest{Dir: dir, ModuleIdent: moduleIdent, DumpFile: dumpFile}
	_, err := h.client.ParseModule(context.Background(), &req)
	return err
}

package main

import (
	"log"
	"net"
	"os"

	"github.com/freeconf/lang/comm"
	"github.com/freeconf/lang/comm/pb"
	"google.golang.org/grpc"
)

const usage = "Usage: %s path-to-socket-file"

func main() {
	if len(os.Args) != 2 {
		log.Fatalf(usage, os.Args[0])
	}
	addr := os.Args[1]
	if _, ferr := os.Stat(addr); ferr == nil {
		if err := os.Remove(addr); err != nil {
			log.Fatal(err)
		}
	}

	l, err := net.Listen("unix", addr)
	if err != nil {
		log.Fatal("listen error:", err)
	}

	s := grpc.NewServer()
	pb.RegisterParserServer(s, &comm.Service{})
	defer l.Close()
	if err := s.Serve(l); err != nil {
		log.Fatalf("grpc server error. %s", err)
	}
}

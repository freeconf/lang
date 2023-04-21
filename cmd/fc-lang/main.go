package main

import (
	"log"
	"os"

	"github.com/freeconf/lang"
)

const usage = `Usage: %s path-to-socket-file path-to-x-socket-file

path-to-socket-file - name of the domain socket file this executable will create
    that will host the gRPC server defined in fc.proto.

path-to-x-socket-file - name of the domain socket file that the program calling
    this executable has hosted the gRPC server defined in fc-x.proto.
`

func main() {
	if len(os.Args) < 3 {
		log.Fatalf(usage, os.Args[0])
	}
	d, err := lang.NewDriver(os.Args[1], os.Args[2])
	chkerr(err)
	chkerr(d.Serve())
}

func chkerr(err error) {
	if err != nil {
		panic(err)
	}
}

package main

import (
	"log"
	"os"

	"github.com/freeconf/lang"
	"github.com/freeconf/lang/pb"
)

const usage = "Usage: %s path-to-socket-file [path-to-x-socket-file]"

func main() {
	if len(os.Args) < 2 {
		log.Fatalf(usage, os.Args[0])
	}
	var c pb.XNodeClient
	if len(os.Args) >= 3 {
		var err error
		c, err = lang.CreateXClient(os.Args[2])
		chkerr(err)
	}
	s, err := lang.NewService(os.Args[1], c)
	chkerr(err)
	chkerr(s.Serve())
}

func chkerr(err error) {
	if err != nil {
		panic(err)
	}
}

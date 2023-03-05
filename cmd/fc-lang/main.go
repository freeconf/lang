package main

import (
	"log"
	"os"

	"github.com/freeconf/lang"
)

const usage = "Usage: %s path-to-socket-file [path-to-x-socket-file]"

func main() {
	if len(os.Args) < 2 {
		log.Fatalf(usage, os.Args[0])
	}
	var gOptionalClientAddr string
	if len(os.Args) >= 3 {
		gOptionalClientAddr = os.Args[2]
	}
	d, err := lang.NewDriver(os.Args[1], gOptionalClientAddr)
	chkerr(err)
	chkerr(d.Serve())
}

func chkerr(err error) {
	if err != nil {
		panic(err)
	}
}

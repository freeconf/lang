package test

import (
	"os"
	"os/exec"

	"github.com/freeconf/lang"
	"github.com/freeconf/yang/fc"
)

type python struct {
	proc   *exec.Cmd
	driver *lang.Driver
}

func (x *python) stop() error {
	fc.Debug.Println("killing python x server")
	return x.proc.Process.Kill()
}

func (x *python) connect(gAddr string, xAddr string) error {
	x.proc = exec.Command("../python/test/harness.py", gAddr, xAddr)
	x.proc.Stdout = os.Stdout
	x.proc.Stderr = os.Stderr
	return x.proc.Start()
}

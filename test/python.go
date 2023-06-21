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
	var err error
	if x.proc != nil {
		fc.Debug.Println("killing python x server")
		err = x.proc.Process.Kill()
		x.proc = nil
	}
	return err
}

func (x *python) connect(gAddr string, xAddr string) error {
	x.proc = exec.Command("../python/tests/harness.py", gAddr, xAddr)
	x.proc.Stdout = os.Stdout
	x.proc.Stderr = os.Stderr
	return x.proc.Start()
}

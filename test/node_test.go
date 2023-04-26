package test

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

var update = flag.Bool("update", false, "update gold files instead of testing against them")

func TestNode(t *testing.T) {

	fc.DebugLog(true)

	tests := []struct {
		Yang   string
		Seed   string
		Expect string
	}{
		{
			Yang:   "basic",
			Seed:   "testdata/seed/basic.json",
			Expect: "testdata/gold/basic.json",
		},
	}
	ypath := source.Dir("testdata/yang")
	h := NewHarness()
	x := &python{}
	fc.RequireEqual(t, nil, h.Connect(x))
	defer x.stop()
	for _, test := range tests {
		m := parser.RequireModule(ypath, test.Yang)
		n := loadSeedData(test.Seed)
		b := node.NewBrowser(m, n)

		dumpFile, err := ioutil.TempFile("", "node-test")
		fc.RequireEqual(t, nil, err)
		fc.RequireEqual(t, nil, h.dump(b.Root(), dumpFile.Name()))
		actual, err := ioutil.ReadFile(dumpFile.Name())
		fc.RequireEqual(t, nil, err)
		fc.Gold(t, *update, actual, test.Expect)
		os.Remove(dumpFile.Name())
	}
}

func loadSeedData(fname string) node.Node {
	raw, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	data := make(map[string]any)
	if err = json.Unmarshal(raw, &data); err != nil {
		panic(err)
	}
	return nodeutil.ReflectChild(data)
}

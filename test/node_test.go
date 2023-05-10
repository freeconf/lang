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

	"github.com/freeconf/lang/pb"
)

var update = flag.Bool("update", false, "update gold files instead of testing against them")

type nodeTestHarness interface {
	createTestCase(tc pb.TestCase, tracefile string) (node.Node, error)
	finalizeTestCase() error
	Close() error
	Connect() error
}

// if you are just running a specific langauge AND remember you change the test data or
// the go implementation, you have to rerun the go and accept the new trace results
// before verify other languages match
var langs = []nodeTestHarness{
	//&golang{},
	NewHarness(&python{}),
}

func TestBasic(t *testing.T) {
	ypath := source.Dir("testdata/yang")
	for _, h := range langs {
		// setup
		fc.RequireEqual(t, nil, h.Connect())
		traceFile := tempFileName()
		n, err := h.createTestCase(pb.TestCase_BASIC, traceFile)
		fc.RequireEqual(t, nil, err)
		m := parser.RequireModule(ypath, "basic")
		b := node.NewBrowser(m, n)

		// test
		err = b.Root().UpsertFrom(readJSON("testdata/seed/basic.json")).LastErr
		fc.AssertEqual(t, nil, err)
		fc.AssertEqual(t, nil, h.finalizeTestCase())
		fc.GoldFile(t, *update, traceFile, "testdata/gold/basic.trace")

		// teardown
		os.Remove(traceFile)
		fc.RequireEqual(t, nil, h.Close())
	}
}

func tempFileName() string {
	dumpFile, err := ioutil.TempFile("", "node-test")
	if err != nil {
		panic(err)
	}
	dumpFile.Close()
	return dumpFile.Name()
}

func TestEcho(t *testing.T) {
	ypath := source.Dir("testdata/yang")
	for _, h := range langs {
		// setup
		fc.RequireEqual(t, nil, h.Connect())
		traceFile := tempFileName()
		n, err := h.createTestCase(pb.TestCase_ECHO, traceFile)
		fc.RequireEqual(t, nil, err)
		m := parser.RequireModule(ypath, "echo")
		b := node.NewBrowser(m, n)

		// test
		sel := b.Root().Find("echo")
		input := nodeutil.ReadJSON(`{
			"f" : 99,
			"g" : {
				"s": "coffee"
			}
		}`)
		output := sel.Action(input)
		fc.AssertEqual(t, nil, output.LastErr)
		fc.AssertEqual(t, nil, h.finalizeTestCase())
		fc.GoldFile(t, *update, traceFile, "testdata/gold/echo.trace")

		// teardown
		os.Remove(traceFile)
		fc.RequireEqual(t, nil, h.Close())
	}
}

func readJSON(fname string) node.Node {
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

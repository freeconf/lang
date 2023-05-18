package test

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"strings"
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

var allLangs = []nodeTestHarness{

	// tip: when debugging in IDE, comment out langs you do not want to debug instead
	// of setting env vars.  Just be sure to restore list when done.

	// newGolang(),

	NewHarness(&python{}),
}

func Langs() []nodeTestHarness {
	langEnv := os.Getenv("FC_LANG")
	if langEnv == "" {
		return allLangs
	}
	var specific []nodeTestHarness
	for _, langId := range strings.Split(langEnv, ",") {
		switch langId {
		case "go":
			specific = append(specific, newGolang())
		case "python":
			specific = append(specific, NewHarness(&python{}))
		}
	}
	return specific
}

func TestBasic(t *testing.T) {
	ypath := source.Dir("testdata/yang")
	for _, h := range Langs() {
		// setup
		fc.RequireEqual(t, nil, h.Connect())
		traceFile := tempFileName()
		n, err := h.createTestCase(pb.TestCase_BASIC, traceFile)
		fc.RequireEqual(t, nil, err)
		m := parser.RequireModule(ypath, "basic")
		b := node.NewBrowser(m, n)

		// test
		fc.AssertEqual(t, nil, b.Root().UpsertFrom(readJSON("testdata/seed/basic.json")))
		fc.AssertEqual(t, nil, h.finalizeTestCase())
		fc.GoldFile(t, *update, traceFile, "testdata/gold/basic.trace")

		// teardown
		os.Remove(traceFile)
		fc.RequireEqual(t, nil, h.Close())
	}
}

func TestEcho(t *testing.T) {
	ypath := source.Dir("testdata/yang")
	for _, h := range Langs() {
		// setup
		fc.RequireEqual(t, nil, h.Connect())
		traceFile := tempFileName()
		n, err := h.createTestCase(pb.TestCase_ECHO, traceFile)
		fc.RequireEqual(t, nil, err)
		m := parser.RequireModule(ypath, "echo")
		b := node.NewBrowser(m, n)

		// test
		sel, err := b.Root().Find("echo")
		fc.RequireEqual(t, nil, err)
		input := nodeutil.ReadJSON(`{
			"f" : 99,
			"g" : {
				"s": "coffee"
			}
		}`)
		output, err := sel.Action(input)
		fc.RequireEqual(t, nil, err)

		echo, err := nodeutil.WritePrettyJSON(output)
		fc.AssertEqual(t, nil, err)
		fc.Gold(t, *update, []byte(echo), "testdata/gold/echo.json")

		fc.AssertEqual(t, nil, h.finalizeTestCase())
		fc.GoldFile(t, *update, traceFile, "testdata/gold/echo.trace")

		// teardown
		os.Remove(traceFile)
		fc.RequireEqual(t, nil, h.Close())
	}
}

func TestNotify(t *testing.T) {
	ypath := source.Dir("testdata/yang")
	m := parser.RequireModule(ypath, "echo")
	for _, h := range Langs() {
		// setup
		fc.RequireEqual(t, nil, h.Connect())
		traceFile := tempFileName()
		n, err := h.createTestCase(pb.TestCase_ECHO, traceFile)
		fc.RequireEqual(t, nil, err)
		b := node.NewBrowser(m, n)

		// test
		sel, err := b.Root().Find("recv")
		fc.RequireEqual(t, nil, err)
		msgs := make(chan string, 1)
		unsub, err := sel.Notifications(func(n node.Notification) {
			msg, err := nodeutil.WritePrettyJSON(n.Event)
			if err != nil {
				panic(err)
			}
			msgs <- msg
		})
		fc.RequireEqual(t, nil, err)

		sendSel, err := b.Root().Find("send")
		fc.RequireEqual(t, nil, err)
		input := nodeutil.ReadJSON(`{
			"f" : 99,
			"g" : {
				"s": "coffee"
			}
		}`)
		_, err = sendSel.Action(input)
		fc.RequireEqual(t, nil, err)
		msg := <-msgs
		fc.Gold(t, *update, []byte(msg), "testdata/gold/send.json")

		unsub()

		fc.AssertEqual(t, nil, h.finalizeTestCase())

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

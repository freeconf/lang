package test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/freeconf/restconf"
	"github.com/freeconf/yang"
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

	parseModule(dir string, module string) (node.Node, error)

	Name() string

	handleCount() int

	Close() error
	Connect() error
}

var goHarness = newGolang(restconf.InternalYPath)
var pythonHarness = NewHarness("python", &python{})

var allHarnesses = []nodeTestHarness{
	//goHarness,
	pythonHarness,
}

func TestBasic(t *testing.T) {
	ypath := source.Dir("testdata/yang")
	m := parser.RequireModule(ypath, "basic")
	for _, h := range Langs() {
		fc.RequireEqual(t, nil, h.Connect())
		t.Run(h.Name(), func(t *testing.T) {
			defer func() {
				fc.RequireEqual(t, nil, h.Close())
			}()
			// setup
			traceFile := tempFileName()
			n, err := h.createTestCase(pb.TestCase_BASIC, traceFile)
			fc.RequireEqual(t, nil, err)
			b := node.NewBrowser(m, n)

			// test
			fc.AssertEqual(t, nil, b.Root().UpsertFrom(readJSON("testdata/seed/basic.json")))
			fc.AssertEqual(t, nil, h.finalizeTestCase())
			fc.GoldFile(t, *update, traceFile, "testdata/gold/basic.trace")

			// teardown
			os.Remove(traceFile)
			fmt.Printf("handles %d\n", h.handleCount())
		})
	}
}

func TestEcho(t *testing.T) {
	ypath := source.Dir("testdata/yang")
	m := parser.RequireModule(ypath, "echo")
	for _, h := range Langs() {
		fc.RequireEqual(t, nil, h.Connect())
		t.Run(h.Name(), func(t *testing.T) {
			defer func() {
				fc.RequireEqual(t, nil, h.Close())
			}()
			// setup
			traceFile := tempFileName()
			n, err := h.createTestCase(pb.TestCase_ECHO, traceFile)
			fc.RequireEqual(t, nil, err)
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
		})
	}
}

func TestNotify(t *testing.T) {
	ypath := source.Dir("testdata/yang")
	m := parser.RequireModule(ypath, "echo")
	for _, h := range Langs() {
		fc.RequireEqual(t, nil, h.Connect())
		t.Run(h.Name(), func(t *testing.T) {
			done := false
			defer func() {
				fc.RequireEqual(t, nil, h.Close())
				done = true
			}()

			// setup
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
					if !done {
						panic(err)
					}
				} else {
					msgs <- msg
				}
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
		})
	}
}

var yangFiles = []string{
	"basic",
	"advanced",
	// "meta",
}

/*
This is very useful test to ensure every facet of a parse yang file is represented
in a given programming language by first parsing a yang in go, then passing those
data structures thru grpc into a given language, then asking for every piece of
those data structures back by having it return every field back.

This can get a little confusing when something breaks however because it relies on
a large part of a language's implementation to work just to get to the part to test.

This will break when Go adds something to yang parser that is also captured in
fc-yang.yang and a language doesn't implement that yet, but that is what this is supposed
to catch to ensure languages have feature parity.
*/
func TestMeta(t *testing.T) {
	dir := "testdata/yang"
	fcYang := parser.RequireModule(yang.InternalYPath, "fc-yang")
	for _, h := range Langs() {
		// setup
		fc.RequireEqual(t, nil, h.Connect())
		t.Run(h.Name(), func(t *testing.T) {
			defer func() {
				fc.RequireEqual(t, nil, h.Close())
			}()

			for _, f := range yangFiles {
				t.Log(f)
				dumpFile, err := os.CreateTemp("", "json")
				fc.RequireEqual(t, nil, err)
				dumper, err := h.parseModule(dir, f)
				fc.RequireEqual(t, nil, err)
				b := node.NewBrowser(fcYang, dumper)
				to := nodeutil.NewJSONWtr(dumpFile)
				to.Pretty = true
				fc.RequireEqual(t, nil, b.Root().UpsertInto(to.Node()))
				goldFile := fmt.Sprintf("testdata/gold/meta/%s.json", f)
				dumpFile.Close()
				fc.GoldFile(t, *update, dumpFile.Name(), goldFile)
				os.Remove(dumpFile.Name())
			}
			fc.AssertEqual(t, nil, h.finalizeTestCase())
		})
	}
}

func TestChoose(t *testing.T) {
	ypath := source.Dir("testdata/yang")
	m := parser.RequireModule(ypath, "advanced")
	for _, h := range Langs() {
		fc.RequireEqual(t, nil, h.Connect())
		t.Run(h.Name(), func(t *testing.T) {
			defer func() {
				fc.RequireEqual(t, nil, h.Close())
			}()
			// setup
			traceFile := tempFileName()
			n, err := h.createTestCase(pb.TestCase_ADVANCED, traceFile)
			fc.RequireEqual(t, nil, err)
			b := node.NewBrowser(m, n)

			// test
			root := b.Root()
			fc.AssertEqual(t, nil, root.UpsertFrom(readJSON("testdata/seed/choose.json")))
			fc.AssertEqual(t, nil, root.UpsertFrom(nodeutil.ReadJSON(`{"two":"dos"}`)))

			two, err := nodeutil.WritePrettyJSON(root)
			fc.AssertEqual(t, nil, err)
			fc.Gold(t, *update, []byte(two), "testdata/gold/choose.json")

			fc.AssertEqual(t, nil, h.finalizeTestCase())
			fc.GoldFile(t, *update, traceFile, "testdata/gold/choose.trace")

			// teardown
			os.Remove(traceFile)
		})
	}
}

func Langs() []nodeTestHarness {
	langEnv := os.Getenv("FC_LANG")
	if langEnv == "" {
		return allHarnesses
	}
	var specific []nodeTestHarness
	for _, langId := range strings.Split(langEnv, ",") {
		switch langId {
		case "go":
			specific = append(specific, goHarness)
		case "python":
			specific = append(specific, pythonHarness)
		}
	}
	return specific
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

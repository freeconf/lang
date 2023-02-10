package codegen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
)

type Meta2Meta struct {
	Definitions []*structDef2
	ByName      map[string]*structDef2
	DataDefs    []*fieldDef2
}

type structDef2 struct {
	Name   string
	Fields []*fieldDef2
}

type fieldDef2 struct {
	Name     string
	Type     string
	Repeated bool
}

func (f *fieldDef2) PyType() string {
	if f.Repeated {
		return f.Type + "[]"
	}
	return f.Type
}

func (f *fieldDef2) PyName() string {
	// avoid python reserved words
	switch f.Name {
	case "def":
		return "ext_def"
	}

	return whisperingSnake(f.Name)
}

func (f *fieldDef2) PyUnpackName() string {
	t := f.PyType()
	if f.Repeated {
		return t[:len(t)-2] + "_array"
	}
	return t
}

func (s *structDef2) IsDataDef() bool {
	switch s.Name {
	case "Container", "List", "Leaf", "LeafList":
		return true
	}
	return false
}

func (s *structDef2) IsMetaDef() bool {
	switch s.Name {
	case "DataDef":
		return false
	}
	return true
}

func title(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func (f *fieldDef2) GoName() string {
	return title(f.Name)
}

func (f *fieldDef2) PyCustomDecoder() string {
	return f.CustomEncoder()
}

func (f *fieldDef2) CustomEncoder() string {
	switch f.Name {
	case "extDef":
		return title(f.Name)
	}
	switch f.Type {
	case "string", "bool", "int32", "int64":
		return ""
	}
	encoder := f.Type
	if f.Repeated {
		encoder = encoder + "List"
	}
	return encoder
}

func ParseMeta2Defs(homeDir string) (Meta2Meta, error) {
	var empty Meta2Meta
	fname := filepath.Join(homeDir, "comm/meta.proto")
	reader, err := os.Open(fname)
	if err != nil {
		return empty, fmt.Errorf("failed to open. %w", err)
	}
	defer reader.Close()

	p := proto.NewParser(reader)
	defs, _ := p.Parse()
	if err != nil {
		return empty, fmt.Errorf("failed to parse. %w", err)
	}
	w := &walker{
		Meta: Meta2Meta{
			ByName: make(map[string]*structDef2),
		},
	}
	proto.Walk(defs, w.handle)
	return w.Meta, nil
}

type walker struct {
	Meta    Meta2Meta
	current *structDef2
}

func (w *walker) handle(v proto.Visitee) {
	switch x := v.(type) {
	case *proto.NormalField:
		w.field(x)
	case *proto.Message:
		w.message(x)
	case *proto.OneOfField:
		w.oneOfField(x)
	}
}

func (w *walker) oneOfField(pf *proto.OneOfField) {
	f := &fieldDef2{
		Name: pf.Name,
		Type: pf.Type,
	}
	w.Meta.DataDefs = append(w.Meta.DataDefs, f)
}

func (w *walker) message(msg *proto.Message) {
	w.current = &structDef2{
		Name: msg.Name,
	}
	w.Meta.Definitions = append(w.Meta.Definitions, w.current)
	w.Meta.ByName[msg.Name] = w.current
}

func (w *walker) field(pf *proto.NormalField) {
	f := &fieldDef2{
		Name:     pf.Name,
		Type:     pf.Type,
		Repeated: pf.Repeated,
	}
	w.current.Fields = append(w.current.Fields, f)
}

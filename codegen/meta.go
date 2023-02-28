package codegen

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/emicklei/proto"
)

type MetaMeta struct {
	Definitions []*structDef
	ByName      map[string]*structDef
	DataDefs    []*fieldDef
	Enums       map[string]*enumDef
}

type structDef struct {
	Name   string
	Fields []*fieldDef
}

type fieldDef struct {
	Name     string
	Type     string
	Repeated bool
}

func (f *fieldDef) PyType() string {
	if f.Repeated {
		return f.Type + "[]"
	}
	return f.Type
}

func (f *fieldDef) PyName() string {
	// avoid python reserved words
	switch f.Name {
	case "def":
		return "ext_def"
	}

	return whisperingSnake(f.Name)
}

func (f *fieldDef) PyUnpackName() string {
	t := f.PyType()
	if f.Repeated {
		return t[:len(t)-2] + "_array"
	}
	return t
}

func (s *structDef) IsDataDef() bool {
	switch s.Name {
	case "Container", "List", "Leaf", "LeafList":
		return true
	}
	return false
}

func (s *structDef) IsMetaDef() bool {
	switch s.Name {
	case "DataDef":
		return false
	}
	return true
}

func (f *fieldDef) GoName() string {
	return title(f.Name)
}

func (f *fieldDef) PyCustomDecoder() string {
	return f.CustomEncoder()
}

func (f *fieldDef) CustomEncoder() string {
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

func ParseMetaDefs(homeDir string) (MetaMeta, error) {
	var empty MetaMeta
	fname := filepath.Join(homeDir, "proto/meta.proto")
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
		Meta: MetaMeta{
			ByName: make(map[string]*structDef),
		},
	}
	proto.Walk(defs, w.handle)
	return w.Meta, nil
}

type enumEntry struct {
	Ident string
	Value int
}

type enumDef struct {
	Ident   string
	Val     int
	Entries []*enumEntry
}

type walker struct {
	Meta    MetaMeta
	current *structDef
	curEnum *enumDef
}

func (w *walker) handle(v proto.Visitee) {
	switch x := v.(type) {
	case *proto.Enum:
		w.enum(x)
	case *proto.EnumField:
		w.enumField(x)
	case *proto.NormalField:
		w.field(x)
	case *proto.Message:
		w.message(x)
	case *proto.OneOfField:
		w.oneOfField(x)
	}
}

func (w *walker) enum(p *proto.Enum) {
	w.curEnum = &enumDef{Ident: p.Name}
	if w.Meta.Enums == nil {
		w.Meta.Enums = make(map[string]*enumDef)
	}
	w.Meta.Enums[w.curEnum.Ident] = w.curEnum
}

func (w *walker) enumField(pf *proto.EnumField) {
	e := &enumEntry{pf.Name, pf.Integer}
	w.curEnum.Entries = append(w.curEnum.Entries, e)
}

func (w *walker) oneOfField(pf *proto.OneOfField) {
	f := &fieldDef{
		Name: pf.Name,
		Type: pf.Type,
	}
	w.Meta.DataDefs = append(w.Meta.DataDefs, f)
}

func (w *walker) message(msg *proto.Message) {
	w.current = &structDef{
		Name: msg.Name,
	}
	w.Meta.Definitions = append(w.Meta.Definitions, w.current)
	w.Meta.ByName[msg.Name] = w.current
}

func (w *walker) field(pf *proto.NormalField) {
	f := &fieldDef{
		Name:     pf.Name,
		Type:     pf.Type,
		Repeated: pf.Repeated,
	}
	w.current.Fields = append(w.current.Fields, f)
}

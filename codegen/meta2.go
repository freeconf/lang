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
	w := new(walker)
	proto.Walk(defs, proto.WithMessage(w.message), proto.WithNormalField(w.field))
	return w.Meta, nil
}

type walker struct {
	Meta    Meta2Meta
	current *structDef2
}

func (w *walker) message(msg *proto.Message) {
	w.current = &structDef2{
		Name: msg.Name,
	}
	w.Meta.Definitions = append(w.Meta.Definitions, w.current)
}

func (w *walker) field(pf *proto.NormalField) {
	f := &fieldDef2{
		Name:     pf.Name,
		Type:     pf.Type,
		Repeated: pf.Repeated,
	}
	w.current.Fields = append(w.current.Fields, f)
}

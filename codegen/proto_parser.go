package codegen

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/emicklei/proto"
)

func ParseProtos(homeDir string) (Vars, error) {
	vars := Vars{
		MessagesByName: make(map[string]*messageDef),
	}
	w := &walker{
		vars: &vars,
	}
	protos := []string{"proto/fc/pb/meta.proto", "proto/fc/pb/val.proto"}
	for _, pname := range protos {
		w.Proto = pname
		fname := filepath.Join(homeDir, pname)
		reader, err := os.Open(fname)
		if err != nil {
			return vars, fmt.Errorf("failed to open. %w", err)
		}
		defer reader.Close()

		p := proto.NewParser(reader)
		defs, _ := p.Parse()
		if err != nil {
			return vars, fmt.Errorf("failed to parse. %w", err)
		}
		proto.Walk(defs, w.handle)
	}

	fmtDefs := w.vars.AllEnums["Format"].Entries
	w.vars.ValEnums = make([]*valEnumEntry, len(fmtDefs)-1)
	for i, def := range fmtDefs[1:] {
		w.vars.ValEnums[i] = &valEnumEntry{
			Ident: def.Ident,
			Value: def.Value,
		}
	}

	for _, msg := range vars.AllMessages {
		if msg.Proto == "proto/fc/pb/meta.proto" {
			mdef := &metaDef{
				Message: msg,
			}
			vars.MetaDefs = append(vars.MetaDefs, mdef)
		}
	}

	return vars, nil
}

type walker struct {
	vars         *Vars
	Proto        string
	current      *messageDef
	currentOneOf string
	curEnum      *enumDef
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
	case *proto.Oneof:
		w.oneOf(x)
	case *proto.OneOfField:
		w.oneOfField(x)
	}
}

func (w *walker) enum(p *proto.Enum) {
	w.curEnum = &enumDef{Ident: p.Name}
	if w.vars.AllEnums == nil {
		w.vars.AllEnums = make(map[string]*enumDef)
	}
	w.vars.AllEnums[w.curEnum.Ident] = w.curEnum
}

func (w *walker) enumField(pf *proto.EnumField) {
	e := &enumEntry{pf.Name, pf.Integer}
	w.curEnum.Entries = append(w.curEnum.Entries, e)
}

func (w *walker) oneOf(pf *proto.Oneof) {
	w.currentOneOf = pf.Name
	w.current.OneOfs[pf.Name] = make([]*fieldDef, 0)
}

func (w *walker) oneOfField(pf *proto.OneOfField) {
	f := &fieldDef{
		Name: pf.Name,
		Type: pf.Type,
	}
	oneofs := w.current.OneOfs[w.currentOneOf]
	w.current.OneOfs[w.currentOneOf] = append(oneofs, f)
}

func (w *walker) message(msg *proto.Message) {
	w.current = &messageDef{
		Proto:  w.Proto,
		Name:   msg.Name,
		OneOfs: make(map[string][]*fieldDef),
	}
	w.vars.AllMessages = append(w.vars.AllMessages, w.current)
	w.vars.MessagesByName[msg.Name] = w.current
}

func (w *walker) field(pf *proto.NormalField) {
	f := &fieldDef{
		Name:     pf.Name,
		Type:     pf.Type,
		Repeated: pf.Repeated,
	}
	w.current.Fields = append(w.current.Fields, f)
}

type messageDef struct {
	Name   string
	Proto  string
	Fields []*fieldDef
	OneOfs map[string][]*fieldDef
}

func (m *messageDef) Field(name string) *fieldDef {
	for _, f := range m.Fields {
		if f.Name == name {
			return f
		}
	}
	return nil
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

type enumEntry struct {
	Ident string
	Value int
}

type enumDef struct {
	Ident   string
	Val     int
	Entries []*enumEntry
}

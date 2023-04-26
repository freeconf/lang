package codegen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
)

type MetaMeta struct {
	Definitions []*structDef
	ByName      map[string]*structDef
	DataDefs    []*fieldDef
	Enums       map[string]*enumDef
	ValEnums    []*valEnumEntry
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

	fmtDefs := w.Meta.Enums["Format"].Entries
	w.Meta.ValEnums = make([]*valEnumEntry, len(fmtDefs)-1)
	for i, def := range fmtDefs[1:] {
		w.Meta.ValEnums[i] = &valEnumEntry{
			Ident: def.Ident,
			Value: def.Value,
		}
	}

	return w.Meta, nil
}

type enumEntry struct {
	Ident string
	Value int
}

type valEnumEntry struct {
	Ident string
	Value int
}

func (def *valEnumEntry) IsList() bool {
	return strings.HasSuffix(def.Ident, "_LIST")
}

func (def *valEnumEntry) GoFmtId() string {
	s := strings.Replace(camel(def.Ident), "Uint", "UInt", 1)
	return fmt.Sprintf("Fmt%s", s)
}

func (def *valEnumEntry) GoNonListType() string {
	switch def.IdentNonList() {
	case "INT32":
		return "int"
	case "UINT32":
		return "uint"
	case "EMPTY":
		return "val.Value"
	case "BITS", "BINARY":
		return "[]byte"
	case "DECIMAL64":
		return "float64"
	case "IDENTITY_REF":
		return "val.IdentRef"
	case "ENUM":
		return "val.Enum"
	}
	return strings.ToLower(def.IdentNonList())
}

func (def *valEnumEntry) ProtoFullyCustomConvert() bool {
	switch def.IdentNonList() {
	case "EMPTY":
		return true
	}
	return false
}

func (def *valEnumEntry) GoType() string {
	if def.IsList() {
		return "[]" + def.GoNonListType()
	}
	return def.GoNonListType()
}

func (def *valEnumEntry) ValType() string {
	s := camel(def.Ident)
	switch s {
	case "IdentityRef":
		return "IdentRef"
	case "IdentityRefList":
		return "IdentRefList"
	}
	return strings.Replace(s, "Uint", "UInt", 1)
}

func (def *valEnumEntry) ValNonListType() string {
	s := camel(def.IdentNonList())
	switch s {
	case "IdentityRef":
		return "IdentRef"
	}
	return strings.Replace(s, "Uint", "UInt", 1)
}

func (def *valEnumEntry) ProtoType() string {
	t := camel(def.IdentNonList())
	switch t {
	case "IdentityRef":
		return "IdentRef"
	}
	return t
}

func (def *valEnumEntry) PyNonListIdent() string {
	return strings.ToLower(def.IdentNonList())
}

func (def *valEnumEntry) GoToProtoNonListConversionRequired() bool {
	switch def.IdentNonList() {
	case "ENUM", "IDENTITY_REF", "EMPTY":
		return true
	}
	return false
}

func (def *valEnumEntry) GoToProtoConversionRequired() bool {
	return false
}

func (def *valEnumEntry) GoToProtoCast(varname string) string {
	switch def.IdentNonList() {
	case "INT8", "INT16", "INT32":
		return fmt.Sprintf("int32(%s)", varname)
	case "UINT8", "UINT16", "UINT32":
		return fmt.Sprintf("uint32(%s)", varname)
	}
	return varname
}

func (def *valEnumEntry) ProtoToGoConversionRequired() bool {
	switch def.Ident {
	case "EMPTY", "EMPTY_LIST":
		return true
	}
	return false
}

func (def *valEnumEntry) ProtoToGoNonListConversionRequired() bool {
	switch def.IdentNonList() {
	case "EMPTY", "ENUM", "IDENTITY_REF":
		return true
	}
	return false
}

func (def *valEnumEntry) ProtoToGoCast(varname string) string {
	switch def.IdentNonList() {
	case "INT8":
		return fmt.Sprintf("int8(%s)", varname)
	case "INT16":
		return fmt.Sprintf("int16(%s)", varname)
	case "INT32":
		return fmt.Sprintf("int(%s)", varname)
	case "UINT8":
		return fmt.Sprintf("uint8(%s)", varname)
	case "UINT16":
		return fmt.Sprintf("uint16(%s)", varname)
	case "UINT32":
		return fmt.Sprintf("uint(%s)", varname)
	}
	return varname
}

func (def *valEnumEntry) IdentNonList() string {
	if def.IsList() {
		return def.Ident[0 : len(def.Ident)-5]
	}
	return def.Ident
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

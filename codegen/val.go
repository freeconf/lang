package codegen

import (
	"fmt"
	"strings"
)

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

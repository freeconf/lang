package codegen

import (
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

func cTemplateFuncs(funcs template.FuncMap) {
	funcs["cDefType"] = cDefType
	funcs["cFieldType"] = cFieldType
	funcs["cName"] = cName
	funcs["cNameSafe"] = cNameSafe
}

func cDefType(d *structDef) string {
	return "fc_meta_" + whisperingSnake(d.Name)
}

func cName(f *fieldDef) string {
	return whisperingSnake(f.Name)
}

func whisperingSnake(s string) string {
	return strings.ToLower(strcase.ToSnake(s))
}

func cNameSafe(name string) string {
	return strings.ReplaceAll(name, "*", "_ptr")
}

func (f *fieldDef) IsArray() bool {
	return strings.HasPrefix(f.Type, "[]")
}

func (f *fieldDef) BaseType() string {
	if f.IsArray() {
		return f.Type[2:]
	}
	return f.Type
}

func cFieldType(f *fieldDef) string {
	switch f.Type {
	case "string":
		return "char*"
	case "[]string":
		return "char**"
	case "*bool":
		return "bool*"
	case "bool":
		return "bool"
	case "int":
		return "int"
	case "int64":
		return "long"
	case "[]ExtensionDefArg":
		return "fc_meta_ext_def_arg_array"
	}
	switch f.Name {
	case "Definitions":
		return "fc_meta_array"
	case "Extensions":
		return "fc_meta_ext_array"
	case "MetaId":
		return "fc_meta_id"
	}
	ctype := whisperingSnake(f.Type)
	if strings.HasPrefix(f.Type, "[]") {
		ctype = ctype[2:] + "**"
	}
	if strings.HasPrefix(f.Type, "*") {
		ctype = ctype[1:] + "*"
	}

	return "fc_meta_" + ctype
}

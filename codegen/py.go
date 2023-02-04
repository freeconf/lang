package codegen

func (f *fieldDef) PyType() string {
	if f.IsArray() {
		return f.BaseType() + "[]"
	}
	return f.Type
}

func (f *fieldDef) PyName() string {
	return whisperingSnake(f.Name)
}

func (f *fieldDef) PyUnpackName() string {
	t := f.PyType()
	if f.IsArray() {
		return t[:len(t)-2] + "_array"
	}
	return t
}

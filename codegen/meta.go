package codegen

type metaDef struct {
	Message *messageDef
}

func (s *metaDef) Name() string {
	return s.Message.Name
}

func (s *metaDef) Fields() []*fieldDef {
	return s.Message.Fields
}

func (s *metaDef) IsDataDef() bool {
	switch s.Message.Name {
	case "Container", "List", "Leaf", "LeafList", "Choice", "Any":
		return true
	}
	return false
}

func (s *metaDef) IsMetaDef() bool {
	switch s.Message.Name {
	case "DataDef":
		return false
	}
	return true && s.Name() != "MetaPointer"
}

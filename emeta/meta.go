package emeta

type DefType int

const (
	DefTypeModule = iota
	DefTypeContainer
	DefTypeLeaf
)

type Module struct {
	Ident       string
	Description string
	Extensions  []Extension
	DefType     DefType
	Ns          string
	Prefix      string
	Contact     string
	Org         string
	Ref         string
	Ver         string
	Definitions []interface{}
}

type Extension struct {
	Name string
}

type Leaf struct {
	Ident       string
	Description string
	Extensions  []Extension
	DefType     int
	Config      *bool
	Mandatory   *bool
}

type Container struct {
	Ident       string
	Description string
	Extensions  []Extension
	DefType     int
	Config      *bool
	Mandatory   *bool
	Definitions []interface{}
}

type List struct {
	Ident       string
	Description string
	Extensions  []Extension
	DefType     int
	Config      *bool
	Mandatory   *bool
	Definitions []interface{}
}

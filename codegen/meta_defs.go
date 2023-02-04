package codegen

// This follows structure of yang/meta/core.go for the most part but has a few
// differences useful for encoding for other languages and in a structure that
// CBOR encoder can read.

type MetaId int

const (
	MetaIdExtension = MetaId(iota)
	MetaIdExtensionDef
	MetaIdExtensionDefArg
	MetaIdModule
	MetaIdLeaf
	MetaIdLeafList
	MetaIdContainer
	MetaIdList
)

type OptionalBool int

const (
	BoolNotSet = iota
	BoolTrue
	BoolFalse
)

// These are defined in order that when they are generated in C there no circular
// references so be ware order has implications for compilation in other languages

type ExtensionDefArg struct {
	MetaId      MetaId
	Ident       string
	Description string
	Ref         string
	YinElement  bool
	// Extensions []*Extension
}

type ExtensionDef struct {
	MetaId      MetaId
	Ident       string
	Description string
	Ref         string
	Status      int
	Args        []ExtensionDefArg
	// Extensions []*Extension
}

type Extension struct {
	MetaId  MetaId
	Ident   string
	Prefix  string
	Keyword string
	ExtDef  string
	Args    []string
}

type Module struct {
	MetaId      MetaId
	Ident       string
	Description string
	Extensions  []Extension

	Definitions []DataDefinitions

	MemId int64

	Ns      string
	Prefix  string
	Contact string
	Org     string
	Ref     string
	Ver     string
}

type Leaf struct {
	MetaId      MetaId
	Ident       string
	Description string
	Extensions  []Extension

	Config    OptionalBool
	Mandatory OptionalBool
}

type LeafList struct {
	MetaId      MetaId
	Ident       string
	Description string
	Extensions  []Extension

	Config    OptionalBool
	Mandatory OptionalBool
}

type Container struct {
	MetaId      MetaId
	Ident       string
	Description string
	Extensions  []Extension

	Definitions []DataDefinitions

	Config    OptionalBool
	Mandatory OptionalBool
}

type DataDefinitions interface{}

type List struct {
	MetaId      MetaId
	Ident       string
	Description string
	Extensions  []Extension

	Definitions []DataDefinitions

	Config    OptionalBool
	Mandatory OptionalBool
}

package emeta

// This follows structure of yang/meta/core.go for the most part but has a few
// differences useful for encoding for other languages and in a structure that
// CBOR encoder can read.

type EncodingId int

const (
	EncodingIdExtension = EncodingId(iota)
	EncodingIdExtensionDef
	EncodingIdExtensionDefArg
	EncodingIdModule
	EncodingIdLeaf
	EncodingIdLeafList
	EncodingIdContainer
	EncodingIdList
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
	Ident       string
	Description string
	Ref         string
	YinElement  bool
	// Extensions []*Extension
}

type ExtensionDef struct {
	Ident       string
	Description string
	Ref         string
	Status      int
	Args        []ExtensionDefArg
	// Extensions []*Extension
}

type Extension struct {
	Ident   string
	Prefix  string
	Keyword string
	Def     string `reference:"ExensionDef"`
	Args    []string
}

type Module struct {
	EncodingId  EncodingId
	Ident       string
	Description string
	Extensions  []Extension
	Ns          string
	Prefix      string
	Contact     string
	Org         string
	Ref         string
	Ver         string
	Definitions []interface{}
}

type Leaf struct {
	EncodingId  EncodingId
	Ident       string
	Description string
	Extensions  []Extension
	Config      OptionalBool
	Mandatory   OptionalBool
}

type LeafList struct {
	EncodingId  EncodingId
	Ident       string
	Description string
	Extensions  []Extension
	Config      OptionalBool
	Mandatory   OptionalBool
}

type Container struct {
	EncodingId  EncodingId
	Ident       string
	Description string
	Extensions  []Extension
	Config      OptionalBool
	Mandatory   OptionalBool
	Definitions []interface{}
}

type List struct {
	EncodingId  EncodingId
	Ident       string
	Description string
	Extensions  []Extension
	Config      OptionalBool
	Mandatory   OptionalBool
	Definitions []interface{}
}

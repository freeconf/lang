package comm

import (
	"github.com/freeconf/lang/comm/pb"
	"github.com/freeconf/yang/meta"
)

type MetaEncoder struct {
}

func (e *MetaEncoder) encodeOptionalBool(value bool, specified bool) pb.OptionalBool {
    if !specified {
        return pb.OptionalBool_NOT_SPECIFIED
    }
    if value {
        return pb.OptionalBool_TRUE
    }
    return pb.OptionalBool_FALSE
}

func (e *MetaEncoder) encodeExtensionList(parent any, from []*meta.Extension) []*pb.Extension {
    to := make([]*pb.Extension, len(from))    
    for i, x := range from {
        to[i] = e.encodeExtension(parent, x)
    }
    return to
}

func (e *MetaEncoder) encodeDataDefList(parent any, from []meta.Definition) []*pb.DataDef {
    to := make([]*pb.DataDef, len(from))    
    for i, d := range from {
        switch x := d.(type) {






            case *meta.Container:
                to[i] = &pb.DataDef{DefOneof:&pb.DataDef_Container{ Container:e.encodeContainer(parent, x)}}

            case *meta.Leaf:
                to[i] = &pb.DataDef{DefOneof:&pb.DataDef_Leaf{ Leaf:e.encodeLeaf(parent, x)}}

            case *meta.List:
                to[i] = &pb.DataDef{DefOneof:&pb.DataDef_List{ List:e.encodeList(parent, x)}}

            case *meta.LeafList:
                to[i] = &pb.DataDef{DefOneof:&pb.DataDef_LeafList{ LeafList:e.encodeLeafList(parent, x)}}
        
        }
    }
    return to
}

func (e *MetaEncoder) encodeExtensionDefArgList(parent *pb.ExtensionDef, from []*meta.ExtensionDefArg) []*pb.ExtensionDefArg {
    to := make([]*pb.ExtensionDefArg, len(from))    
    for i, f := range from {
        to[i] = e.encodeExtensionDefArg(parent, f)
    }
    return to
}

func (e *MetaEncoder) encodeStatus(parent *pb.ExtensionDef, from meta.Status) pb.Status {
    return pb.Status(from)
}




func (e *MetaEncoder) encodeModule(parent any, from *meta.Module) *pb.Module {
    var def pb.Module
        def.Ident = from.Ident()
        def.Description = from.Description()
        def.Extensions = e.encodeExtensionList(&def, from.Extensions())
        def.Definitions = e.encodeDataDefList(&def, from.DataDefinitions())
        def.Namespace = from.Namespace()
        def.Prefix = from.Prefix()
        def.Contact = from.Contact()
        def.Organization = from.Organization()
        def.Reference = from.Reference()
        def.Version = from.Version()
    return &def
}

func (e *MetaEncoder) encodeExtensionDefArg(parent any, from *meta.ExtensionDefArg) *pb.ExtensionDefArg {
    var def pb.ExtensionDefArg
        def.Ident = from.Ident()
        def.Description = from.Description()
        def.Reference = from.Reference()
        def.YinElement = from.YinElement()
    return &def
}

func (e *MetaEncoder) encodeExtensionDef(parent any, from *meta.ExtensionDef) *pb.ExtensionDef {
    var def pb.ExtensionDef
        def.Ident = from.Ident()
        def.Description = from.Description()
        def.Reference = from.Reference()
        def.Status = e.encodeStatus(&def, from.Status())
        def.Arguments = e.encodeExtensionDefArgList(&def, from.Arguments())
    return &def
}

func (e *MetaEncoder) encodeExtension(parent any, from *meta.Extension) *pb.Extension {
    var def pb.Extension
        def.Ident = from.Ident()
        def.Prefix = from.Prefix()
        def.Keyword = from.Keyword()
        def.Def = from.Def().Ident()
        def.Arguments = from.Arguments()
    return &def
}

func (e *MetaEncoder) encodeContainer(parent any, from *meta.Container) *pb.Container {
    var def pb.Container
        def.Ident = from.Ident()
        def.Description = from.Description()
        def.Extensions = e.encodeExtensionList(&def, from.Extensions())
        def.Definitions = e.encodeDataDefList(&def, from.DataDefinitions())
        def.Config = e.encodeOptionalBool(from.Config(), from.IsConfigSet())
        def.Mandatory = e.encodeOptionalBool(from.Mandatory(), from.IsMandatorySet())
    return &def
}

func (e *MetaEncoder) encodeLeaf(parent any, from *meta.Leaf) *pb.Leaf {
    var def pb.Leaf
        def.Ident = from.Ident()
        def.Description = from.Description()
        def.Extensions = e.encodeExtensionList(&def, from.Extensions())
        def.Config = e.encodeOptionalBool(from.Config(), from.IsConfigSet())
        def.Mandatory = e.encodeOptionalBool(from.Mandatory(), from.IsMandatorySet())
    return &def
}

func (e *MetaEncoder) encodeList(parent any, from *meta.List) *pb.List {
    var def pb.List
        def.Ident = from.Ident()
        def.Description = from.Description()
        def.Extensions = e.encodeExtensionList(&def, from.Extensions())
        def.Definitions = e.encodeDataDefList(&def, from.DataDefinitions())
        def.Config = e.encodeOptionalBool(from.Config(), from.IsConfigSet())
        def.Mandatory = e.encodeOptionalBool(from.Mandatory(), from.IsMandatorySet())
    return &def
}

func (e *MetaEncoder) encodeLeafList(parent any, from *meta.LeafList) *pb.LeafList {
    var def pb.LeafList
        def.Ident = from.Ident()
        def.Description = from.Description()
        def.Extensions = e.encodeExtensionList(&def, from.Extensions())
        def.Config = e.encodeOptionalBool(from.Config(), from.IsConfigSet())
        def.Mandatory = e.encodeOptionalBool(from.Mandatory(), from.IsMandatorySet())
    return &def
}


func (e *MetaEncoder) Encode(from *meta.Module) *pb.Module {
    return e.encodeModule(nil, from)
}

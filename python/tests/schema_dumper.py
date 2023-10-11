
from freeconf import nodeutil, meta, val

# The only purpose of this is to catch missing fields and data structures in freeconf.meta
# package. It has no purpose outside that because you can use Go's nodeutil.Schema instead
# which should match the functionality here.  If there ever is a need to have that 
# implemented in python, this would be a great start.

class Dumper:

    def node(self, m):

        def choose(n, sel, choice):
            if choice.ident == "body-stmt":
                if isinstance(n.object, meta.Module):
                    return choice.cases['module']
                if isinstance(n.object, meta.Container):
                    return choice.cases['container']
                if isinstance(n.object, meta.List):
                    return choice.cases['list']
                if isinstance(n.object, meta.Choice):
                    return choice.cases['choice']
                if isinstance(n.object, meta.ChoiceCase):
                    return choice.cases['case']
                if isinstance(n.object, meta.Leaf):
                    return choice.cases['leaf']
                if isinstance(n.object, meta.LeafList):
                    return choice.cases['leaf-list']
                if isinstance(n.object, meta.Any):
                    return choice.cases['anyxml']
                if isinstance(n.object, SchemaPointer):
                    return choice.cases['pointer']
            return n.do_choose(sel, choice)

        def child(n, r):
            if r.meta.ident in ['module', 'container', 'list', 'choice', 'case', 'leaf-list', 'leaf', 'anyxml', 'pointer']:
                return n.new(n.object)
            elif r.meta.ident == 'dataDef':
                if self.has_recursive_child(n.object):
                    copy = []
                    for ddef in n.object.definitions:
                        if ddef.parent != n.object:
                            copy.append(SchemaPointer(n.object, ddef))
                        else:
                            copy.append(ddef)                    

            return n.do_child(r)
        
        def field(n, r, v):
            if r.meta.ident == "status":
                if isinstance(n.object, meta.Module) or isinstance(n.object, meta.ExtensionDefArg):
                    return None
            elif r.meta.ident == "when":
                if n.object.when != None:
                    return val.Val(n.object.when.expression)
            #  for default, add to meta/protos
            return n.do_field(r, v)
        
        def options(n, m, opts):
            aliases = {
                "notify": "notifications",
                "dataDef": "definitions",
                "identity": "identities",
                "default": "default_val",
                "defaults" : "default_vals",
                "enumeration" : "enums"
            }
            opts.ident = aliases.get(m.ident, None)
            return opts
                    
        opts = nodeutil.NodeOptions(
            try_plural_on_lists=True
        )
        return nodeutil.Node(m, 
            options=opts,
            on_child=child,
            on_choose=choose,
            on_field=field, 
            on_options=options)
    

    def has_recursive_child(self, m):
        for ddef in m.definitions:
            if ddef.parent != m:
                return True
        return False


class SchemaPointer:

    def __init__(self, parent, delegate):
        self.parent = parent
        self.delegate = delegate

import pb.fc_g_pb2
import pb.fc_g_pb2_grpc
import pb.fc_x_pb2
import pb.fc_x_pb2_grpc
import fc.handles
import fc.meta
import fc.val

class Selection():

    def __init__(self, node_service, node, path, parent, browser):
        self.node_service = node_service
        self.parent = parent
        self.hnd = 0
        self.node = node
        self.path = path
        self.browser = browser

    @classmethod
    def new_split(cls, node_service, browser, path, node):
        return Selection(node_service, node, path, None, browser)


    @classmethod
    def new_root(cls, node_service, browser, hnd=None):
        path = fc.meta.Path(None, browser.module)
        sel = Selection(node_service, browser.node, path, None, browser)
        sel.hnd = fc.handles.put(sel, hnd)
        return sel


    def new_select(self, meta, node, hnd=None):
        path = fc.meta.Path(self.path, meta)
        sel = Selection(self.node_service, node, path, self, self.browser)
        sel.hnd = fc.handles.put(sel, hnd)
        return sel


    def action(self, n=0):
        nodeHnd = self.node_service.lazy_node_hnd(n)
        req = pb.fc_g_pb2.ActionRequest(selHnd=self.hnd, inputNodeHnd=nodeHnd)
        resp = self.node_service.stub.Action(req)
        outputSel = None
        if resp.outputSelHnd:
            outputNode = self.node_service.lazy_node_hnd(self, resp.outputHnd)
            outputSel = self.new_select(self.path.meta.output, outputNode, resp.outputSelHnd)
        return outputSel


    def upsert_from(self, n):
        req = pb.fc_g_pb2.UpsertFromRequest(selHnd=self.hnd)        
        req.nodeHnd = self.node_service.lazy_node_hnd(n)
        print(f'upsert_from sel_hnd={req.selHnd}, sel.node.hnd={self.node.hnd}, node_hnd={n.hnd}')
        self.node_service.stub.UpsertFrom(req)


    def find(self, path):
        req = pb.fc_g_pb2.FindRequest(selHnd=self.hnd, path=path)
        resp = self.node_service.stub.Find(req)
        found = self.node_service.resolve_sel(resp.selHnd)
        return found




class Browser():

    def __init__(self, node_service, browser_hnd, module, node=None):
        self.node_service = node_service
        self.hnd = browser_hnd
        self.module = module
        self.node = node


    def root(self):
        g_req = pb.fc_g_pb2.BrowserRootRequest(browserHnd=self.hnd)
        resp = self.node_service.stub.BrowserRoot(g_req)
        return self.node_service.resolve_sel(resp.selHnd)


class NodeService():

    def __init__(self, driver):
        self.driver = driver
        self.stub = pb.fc_g_pb2_grpc.NodeStub(driver.channel)

    def new_browser(self, m, n):
        node_hnd = self.lazy_node_hnd(n)
        req = pb.fc_g_pb2.NewBrowserRequest(moduleHnd=m.hnd, nodeHnd=node_hnd)
        resp = self.stub.NewBrowser(req)
        b = Browser(self, resp.browserHnd, m, n)
        b.hnd = fc.handles.put(b, resp.browserHnd)
        return b

    def lazy_node_hnd(self, n):
        if not n:
            return 0
        if isinstance(n, int):
            return fc.handles.get(n)
        if not n.hnd:
            resp = self.stub.NewNode(pb.fc_g_pb2.NewNodeRequest())
            n.hnd = resp.nodeHnd
            fc.handles.put(n, n.hnd)
        return n.hnd


    def resolve_sel(self, sel_hnd):
        try:
            return fc.handles.get(sel_hnd)
        except KeyError:
            req = pb.fc_g_pb2.GetSelectionRequest(selHnd=sel_hnd)
            resp = self.stub.GetSelection(req)
            if resp.parentHnd:
                parent = self.resolve_sel(resp.parentHnd)
                meta = fc.meta.require_def(parent.path.meta, resp.metaIdent)
                node = self.lazy_node_hnd(resp.nodeHnd)
                sel = parent.new_select(meta, node, sel_hnd)
            else:
                browser = self.resolve_browser(resp.browserHnd)
                sel = Selection.new_root(self, browser, sel_hnd)
            return sel


    def resolve_browser(self, browser_hnd):
        try:
            return fc.handles.get(browser_hnd)
        except KeyError:
            req = pb.fc_g_pb2.GetBrowserRequest(browserHnd=browser_hnd)
            resp = self.stub.GetBrowser(req)
            module = self.resolve_module(resp.moduleHnd)
            return Browser(self, browser_hnd, module)


    def resolve_module(self, module_hnd):
        try:
            fc.handles.get(module_hnd)
        except KeyError:
            req = pb.fc_g_pb2.GetModuleRequest(moduleHnd=module_hnd)
            resp = self.stub.GetModule(req)
            m = fc.meta_decoder.Decoder().decode(resp.module)
            m.hnd = module_hnd
            fc.handles.put(m, module_hnd)
            return m

class ChildRequest():

    def __init__(self, sel, meta, new, delete):
        self.sel = sel
        self.meta = meta
        self.new = new
        self.delete = delete

class ActionRequest():

    def __init__(self, sel, meta, input):
        self.sel = sel
        self.meta = meta
        self.input = input        


class FieldRequest():

    def __init__(self, sel, meta, write, clear):
        self.sel = sel
        self.meta = meta
        self.write = write
        self.clear = clear


class XNodeServicer(pb.fc_x_pb2_grpc.XNodeServicer):
    """Bridge between python node navigation and go node navigation"""

    def __init__(self):
        self.node_service = None

    def Child(self, g_req, context):
        sel = self.node_service.resolve_sel(g_req.selHnd)
        meta = fc.meta.require_def(sel.path.meta, g_req.metaIdent)
        req = ChildRequest(sel, meta, g_req.new, g_req.delete)
        child = sel.node.child(req)
        childNodeHnd = self.node_service.lazy_node_hnd(child)
        return pb.fc_x_pb2.ChildResponse(nodeHnd=childNodeHnd)


    def Field(self, g_req, context):
        sel = self.node_service.resolve_sel(g_req.selHnd)
        meta = fc.meta.require_def(sel.path.meta, g_req.metaIdent)
        req = FieldRequest(sel, meta, g_req.write, g_req.clear)
        write_val = None
        if g_req.write:
            write_val = fc.val.proto_decode(g_req.toWrite)
        read_val = sel.node.field(req, write_val)
        if not g_req.write:
            fromRead = fc.val.proto_decode(read_val)
            resp = pb.fc_x_pb2.FieldResponse(fromRead=fromRead)
        else:
            resp = pb.fc_x_pb2.FieldResponse()

        return resp


    def Select(self, g_req, context):
        parent = None
        meta = None
        browser = None
        if g_req.parentSelHnd > 0:
            parent = self.node_service.resolve_sel(g_req.parentSelHnd)
            browser = parent.browser
            meta = fc.meta.get_def(parent.path.meta, g_req.metaIdent)
            node = fc.handles.get(g_req.nodeHnd)
            sel = parent.new_select(meta, node)
        elif g_req.browserHnd > 0:
            browser = fc.handles.get(g_req.browserHnd)
            sel = Selection.new_root(self.node_service, browser)
        else:
            raise Exception('no module or parent selection given')
        resp = pb.fc_x_pb2.SelectResponse(selHnd=sel.hnd)
        return resp

    def Split(self, g_req, context):
        node = fc.handles.get(g_req.nodeHnd)
        browser = None # TODO: could be local browser, could not be
        module = fc.handles.get(g_req.moduleHnd)
        path = fc.meta.new_path(module, g_req.metaPath)
        sel = Selection.new_split(self.node_service, browser, path, node)
        resp = pb.fc_x_pb2.SplitResponse(selHnd=sel.hnd)
        return resp

    def Action(self, g_req, context):
        sel = self.node_service.resolve_sel(g_req.selHnd)
        # TODO: Id this inconsistent?
        # meta = fc.meta.require_def(sel.path.meta, g_req.metaIdent)
        meta = sel.path.meta
        input = None
        if g_req.inputSelHnd:
            input = self.node_service.lazy_node_hnd(g_req.inputSelHnd)
        output = sel.node.action(ActionRequest(sel, meta, input))
        outputNodeHnd = self.node_service.lazy_node_hnd(output)
        return pb.fc_x_pb2.XActionResponse(outputNodeHnd=outputNodeHnd)

import pb.fc_g_pb2
import pb.fc_g_pb2_grpc
import pb.fc_x_pb2
import pb.fc_x_pb2_grpc
import fc.handles
import fc.meta
import fc.val

class Node():

    def __init__(self):
        self.hnd = 0

    def select(self):
        pass

    def child(self):
        pass

    def field(self):
        pass

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
    def new_root(cls, node_service, browser):
        path = fc.meta.Path(None, browser.module)
        sel = Selection(node_service, browser.node, path, None, browser)
        sel.hnd = fc.handles.put(sel)
        return sel


    def new_select(self, meta, node):
        path = fc.meta.Path(self.path, meta)
        sel = Selection(self.node_service, node, path, self, self.browser)
        sel.hnd = fc.handles.put(sel)
        return sel

    def upsert_from(self, n):
        req = pb.fc_g_pb2.UpsertFromRequest(selHnd=self.hnd)
        req.nodeHnd = self.node_service.lazy_node_hnd(n)
        self.node_service.stub.UpsertFrom(req)


class Browser():

    def __init__(self, node_service, module, node, browser_hnd):
        self.node_service = node_service
        self.hnd = browser_hnd
        self.module = module
        self.node = node


    def root(self):
        g_req = pb.fc_g_pb2.BrowserRootRequest(browserHnd=self.hnd)
        resp = self.node_service.stub.BrowserRoot(g_req)
        return fc.handles.get(resp.selHnd)


class NodeService():

    def __init__(self, driver):
        self.driver = driver
        self.stub = pb.fc_g_pb2_grpc.NodeStub(driver.channel)

    def new_browser(self, m, n):
        node_hnd = self.lazy_node_hnd(n)
        req = pb.fc_g_pb2.NewBrowserRequest(moduleHnd=m.hnd, nodeHnd=node_hnd)
        resp = self.stub.NewBrowser(req)
        b = Browser(self, m, n, resp.browserHnd)
        b.hnd = fc.handles.put(b, resp.browserHnd)
        return b

    def lazy_node_hnd(self, n):
        if isinstance(n, int):
            return n
        if not n.hnd:
            resp = self.stub.NewNode(pb.fc_g_pb2.NewNodeRequest())
            n.hnd = resp.nodeHnd
            fc.handles.put(n, n.hnd)
        return n.hnd

class ChildRequest():

    def __init__(self, sel, meta, new, delete):
        self.sel = sel
        self.meta = meta
        self.new = new
        self.delete = delete


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
        sel = fc.handles.get(g_req.selHnd)
        meta = fc.meta.require_def(sel.path.meta, g_req.metaIdent)
        req = ChildRequest(sel, meta, g_req.new, g_req.delete)
        child = sel.node.child(req)
        g_resp = pb.fc_x_pb2.ChildResponse()
        if child:
            child.hnd = self.node_service.lazy_node_hnd(child)
            g_resp.nodeHnd = child.hnd
        return g_resp


    def Field(self, g_req, context):
        sel = fc.handles.get(g_req.selHnd)
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
            parent = fc.handles.get(g_req.parentSelHnd)
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

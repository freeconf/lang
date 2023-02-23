import pb.fc_g_pb2
import pb.fc_g_pb2_grpc
import pb.fc_x_pb2
import pb.fc_x_pb2_grpc
import fc.handles
import fc.meta
import fc.val

class Node():
    __stubs__ = [
        "x_node_hnd"
    ]

    def select(self):
        pass

    def child(self):
        pass

    def field(self):
        pass


class Path():

    def __init__(self, parent, meta):
        self.parent = parent
        self.meta = meta


class Selection():

    def __init__(self, node_service, parent, browser, g_sel_hnd, meta, node):
        self.node_service = node_service
        self.parent = parent
        self.browser = browser
        self.g_sel_hnd = g_sel_hnd
        if parent:
            self.path = Path(parent.path, meta)
        else:
            self.path = Path(None, meta)
        self.node = node


    def sel_child(self, g_sel_hnd, meta, node):
        return Selection(self.node_service, self, self.browser, g_sel_hnd, meta, node)


    def upsert_from(self, n):
        req = pb.fc_g_pb2.UpsertFromRequest(gSelHnd=self.g_sel_hnd)
        if isinstance(n, int):
            req.gNodeHnd = n
        else:
            req.xNodeHnd = n.x_node_hnd
        self.node_service.stub.UpsertFrom(req)
        

class Browser():

    def __init__(self, node_service, module, node):
        self.node_service = node_service
        self.g_browser_hnd = None
        self.x_browser_hnd = None
        self.module = module
        self.node = node


    def root(self):
        g_req = pb.fc_g_pb2.BrowserRootRequest(gBrowserHnd=self.g_browser_hnd, xBrowserHnd=self.x_browser_hnd)
        g_resp = self.node_service.stub.BrowserRoot(g_req)
        sel = Selection(self.node_service, None, self, g_resp.gSelHnd, self.module, self.node)
        sel.handle = fc.handles.put(sel)
        return sel


class NodeService():

    def __init__(self, driver):
        self.driver = driver
        self.stub = pb.fc_g_pb2_grpc.NodeStub(driver.channel)


    def new_browser(self, m, n):
        x_node_hnd = fc.handles.put(n)
        b = Browser(self, m, n)
        b.x_browser_hnd = fc.handles.put(b)
        req = pb.fc_g_pb2.NewBrowserRequest(gModuleHnd=m.handle, xNodeHnd=x_node_hnd, xBrowserHnd=b.x_browser_hnd)
        resp = self.stub.NewBrowser(req)
        b.g_browser_hnd = resp.gBrowserHnd
        return b


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

    def __init__(self, *args, **kwargs):        
        pass

    def Child(self, g_req, context):
        sel = fc.handles.get(g_req.xSelHnd)
        meta = fc.meta.get_def(sel.path.meta, g_req.metaIdent)
        req = ChildRequest(sel, meta, g_req.new, g_req.delete)
        child = sel.node.child(req)
        g_resp = pb.fc_x_pb2.ChildResponse()
        if child:
            child.x_node_hand = fc.handles.put(child)
            g_resp.xNodeHnd = child.x_node_hand
        return g_resp


    def Field(self, g_req, context):
        sel = fc.handles.get(g_req.xSelHnd)
        meta = fc.meta.get_def(sel.path.meta, req.meta)
        req = FieldRequest(sel, meta, req.new, req.delete)
        write_val = None
        if req.write:
            write_val = fc.val.proto_decode(req.toWrite)
        read_val = sel.node.field(req, write_val)
        resp = pb.fc_x_pb2.FieldResponse()
        if not req.write:
            resp.fromRead = fc.val.proto_decode(read_val)
        return resp


    def Select(self, g_req, context):
        parent = None
        meta = None
        browser = None
        if g_req.xSelHnd > 0:
            parent = fc.handles.get(g_req.xSelHnd)
            browser = parent.browser
            meta = fc.meta.get_def(parent.path.meta, g_req.metaIdent)
        elif g_req.xBrowserHnd > 0:
            browser = fc.handles.get(g_req.xBrowserHnd)
            meta = browser.module
        else:
            raise Exception('no module or parent selection given')
        node = fc.handles.get(g_req.xNodeHnd)
        sel = Selection(self, parent, browser, g_req.gSelHnd, meta, node)
        resp = pb.fc_x_pb2.SelectResponse(xSelHnd=fc.handles.put(sel))
        return resp


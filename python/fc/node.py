import queue
import grpc
import threading
import pb.fc_g_pb2
import pb.fc_g_pb2_grpc
import pb.fc_x_pb2
import pb.fc_x_pb2_grpc
import fc.handles
import fc.meta
import fc.val
import fc.parser
import traceback

class Selection():

    def __init__(self, driver, hnd_id, node, path, parent, browser):
        self.driver = driver
        self.hnd = fc.handles.Handle(driver, hnd_id, self)
        self.parent = parent
        self.node = node
        self.path = path
        self.browser = browser

    @classmethod
    def resolve(cls, driver, hnd_id):
        try:
            return fc.handles.Handle.require(driver, hnd_id)
        except KeyError:
            req = pb.fc_g_pb2.GetSelectionRequest(selHnd=hnd_id)
            resp = driver.g_nodes.GetSelection(req)
            node = resolve_node(driver, resp.nodeHnd)
            if resp.parentHnd:
                parent = Selection.resolve(driver, resp.parentHnd)
                if isinstance(parent.path.meta, fc.meta.Notification):
                    meta = parent.path.meta
                    path = fc.meta.Path(parent.path, meta)
                else:
                    meta = fc.meta.require_def(parent.path.meta, resp.metaIdent)
                    path = fc.meta.Path(parent.path, meta)
                sel = Selection(driver, hnd_id, node, path, parent, parent.browser) 
            else:
                browser = fc.node.Browser.resolve(driver, resp.browserHnd)
                path = fc.meta.Path(None, browser.module)
                sel = Selection(driver, hnd_id, node, path, None, browser)
            return sel

    def action(self, inputNode=0):
        inputNodeHnd = 0        
        if inputNode != 0:
            inputNodeHnd = resolve_node(self.driver, inputNode).hnd.id
        req = pb.fc_g_pb2.ActionRequest(selHnd=self.hnd.id, inputNodeHnd=inputNodeHnd)
        resp = self.driver.g_nodes.Action(req)
        outputSel = None
        if resp.outputSelHnd:
            outputSel = Selection.resolve(self.driver, resp.outputSelHnd)
        return outputSel
    
    def notification(self, callback):
        req = pb.fc_g_pb2.NotificationRequest(selHnd=self.hnd.id)
        stream = self.driver.g_nodes.Notification(req)
        def rdr():
            try:
                for resp in stream:
                    if resp == None:
                        return
                    msg = Selection.resolve(self.driver, resp.selHnd)
                    callback(msg)
            except grpc.RpcError as gerr:
                if not gerr.cancelled():
                    print(f'grpc err. {gerr}')
            except Exception as e:
                print(f'got error in callback delivering msg: {type(e)} {e}')

        t = threading.Thread(target=rdr)
        t.start()
        def closer():
            stream.cancel()
            t.join()
        return closer

    def upsert_from(self, n):
        n = resolve_node(self.driver, n)
        req = pb.fc_g_pb2.UpsertFromRequest(selHnd=self.hnd.id, nodeHnd=n.hnd.id)        
        self.driver.g_nodes.UpsertFrom(req)

    def find(self, path):
        req = pb.fc_g_pb2.FindRequest(selHnd=self.hnd.id, path=path)
        resp = self.driver.g_nodes.Find(req)
        return Selection.resolve(self.driver, resp.selHnd)


def resolve_node(driver, n):
    if not n:
        # nil node
        return 0
    if isinstance(n, int):
        # cached node
        return fc.handles.Handle.require(driver, n)
    if not n.hnd:
        # unregistered local node about to be registered with go
        resp = driver.g_nodes.NewNode(pb.fc_g_pb2.NewNodeRequest())
        n.hnd = fc.handles.Handle(driver, resp.nodeHnd, n)
    return n


class Browser():

    def __init__(self, driver, module, node=None, node_src=None, hnd_id=None):
        self.driver = driver
        if not hnd_id:
            req = pb.fc_g_pb2.NewBrowserRequest(moduleHnd=module.hnd.id)
            resp = self.driver.g_nodes.NewBrowser(req)
            self.hnd = fc.handles.Handle(driver, resp.browserHnd, self)
        else:
            self.hnd = fc.handles.Handle(driver, hnd_id, self)
        self.module = module
        self.node_src = node_src
        self.node_obj = node

    def node(self):
        if self.node_src:
            return self.node_src()
        return self.node_obj

    def root(self):
        g_req = pb.fc_g_pb2.BrowserRootRequest(browserHnd=self.hnd.id)
        resp = self.driver.g_nodes.BrowserRoot(g_req)
        return Selection.resolve(self.driver, resp.selHnd)

    @classmethod
    def resolve(cls, driver, hnd_id):
        try:
            return fc.handles.Handle.require(driver, hnd_id)
        except KeyError:
            req = pb.fc_g_pb2.GetBrowserRequest(browserHnd=hnd_id)
            resp = driver.g_nodes.GetBrowser(req)
            module = fc.parser.Parser.resolve_module(driver, resp.moduleHnd)
            return Browser(driver, module, hnd_id=hnd_id)


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

class NotificationRequest():

    def __init__(self, sel, meta, queue):
        self.sel = sel
        self.meta = meta
        self.queue = queue

    def send(self, node):
        self.queue.put(node)

class NotificationSession():

    def __init__(self, q, closer):
        self.q = q
        self.closer = closer
        self.hnd = None

    @classmethod
    def resolve(cls, driver, hnd_id):
        return fc.handles.Handle.require(driver, hnd_id)


class XNodeServicer(pb.fc_x_pb2_grpc.XNodeServicer):
    """Bridge between python node navigation and go node navigation"""

    def __init__(self, driver):
        self.driver = driver

    def Child(self, g_req, context):
        sel = Selection.resolve(self.driver, g_req.selHnd)
        meta = fc.meta.require_def(sel.path.meta, g_req.metaIdent)
        req = ChildRequest(sel, meta, g_req.new, g_req.delete)
        child = sel.node.child(req)
        child = resolve_node(self.driver, child)
        return pb.fc_x_pb2.ChildResponse(nodeHnd=child.hnd.id)

    def Field(self, g_req, context):
        try:
            sel = Selection.resolve(self.driver, g_req.selHnd)
            meta = fc.meta.require_def(sel.path.meta, g_req.metaIdent)
            req = FieldRequest(sel, meta, g_req.write, g_req.clear)
            write_val = None
            if g_req.write:
                write_val = fc.val.proto_decode(g_req.toWrite)
            read_val = sel.node.field(req, write_val)
            if not g_req.write:
                fromRead = fc.val.proto_encode(read_val)
                resp = pb.fc_x_pb2.FieldResponse(fromRead=fromRead)
            else:
                resp = pb.fc_x_pb2.FieldResponse()
            return resp
        except Exception as error:
            print(traceback.format_exc())
            raise error

    def Select(self, g_req, context):
        # TODO
        pass

    def Action(self, g_req, context):
        sel = Selection.resolve(self.driver, g_req.selHnd)
        # TODO: Id this inconsistent?
        # meta = fc.meta.require_def(sel.path.meta, g_req.metaIdent)
        meta = sel.path.meta
        input = None
        if g_req.inputSelHnd:
            input = resolve_node(self.driver, g_req.inputSelHnd)
        output = sel.node.action(ActionRequest(sel, meta, input))
        outputNodeHnd = None
        if output:
            outputNodeHnd = resolve_node(self.driver, g_req.outputNodeHnd).hnd.id
        return pb.fc_x_pb2.XActionResponse(outputNodeHnd=outputNodeHnd)


    def NodeSource(self, g_req, context):
        browser = fc.node.Browser.resolve(self.driver, g_req.browserHnd)
        n = resolve_node(self.driver, browser.node())
        return pb.fc_x_pb2.NodeSourceResponse(nodeHnd=n.hnd.id)


    def Notification(self, g_req, context):
        q = queue.Queue()
        sel = Selection.resolve(self.driver, g_req.selHnd)
        meta = sel.path.meta
        closer = sel.node.notification(NotificationRequest(sel, meta, q))
        try:
            while True:
                node = q.get()
                if node == None:
                    break
                n = resolve_node(self.driver, node)
                yield pb.fc_x_pb2.XNotificationResponse(nodeHnd=n.hnd.id)
                q.task_done()
        finally:
            closer()
        return None

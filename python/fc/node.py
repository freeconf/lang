import queue
import grpc
import threading
import pb.fc_pb2
import pb.common_pb2
import pb.fc_pb2_grpc
import pb.fc_x_pb2
import pb.fc_x_pb2_grpc
import fc.handles
import fc.meta
import fc.val
import fc.parser
import traceback
import logging

class Selection():

    def __init__(self, driver, hnd_id, node, path, parent, browser, inside_list=False):
        self.driver = driver
        self.hnd = fc.handles.Handle(driver, hnd_id, self)
        self.parent = parent
        self.node = node
        self.path = path
        self.browser = browser
        self.inside_list = inside_list

    @classmethod
    def resolve(cls, driver, hnd_id):
        sel = fc.handles.Handle.lookup(driver, hnd_id)
        if sel == None:
            req = pb.fc_pb2.GetSelectionRequest(selHnd=hnd_id)
            resp = driver.g_nodes.GetSelection(req)
            print(f'Selection.resolve hnd={hnd_id}, resp.parentHnd={resp.parentHnd}')
            node = resolve_node(driver, resp.nodeHnd)
            if resp.parentHnd:
                # recursive
                print(f'resolve w/parent hnd={hnd_id} path.type={resp.path.type}, input?={resp.path.type == pb.common_pb2.RPC_INPUT}')
                parent = Selection.resolve(driver, resp.parentHnd)
                if resp.path.type == pb.common_pb2.DATA_DEF:
                    ddef_meta = fc.meta.get_def(parent.path.meta, resp.path.metaIdent)
                    path = fc.meta.Path(parent.path, ddef_meta)
                elif resp.path.type == pb.common_pb2.LIST_ITEM:
                    key = None
                    if resp.path.key != None and len(resp.path.key) > 0:
                        key = [fc.val.proto_decode(v) for v in resp.path.key]
                    meta = parent.path.meta
                    path = fc.meta.Path(parent.path, meta, key=key)
                elif resp.path.type == pb.common_pb2.RPC:
                    rpc_meta = fc.meta.get_rpc(parent.path.meta, resp.path.metaIdent)
                    path = fc.meta.Path(parent.path, rpc_meta)
                elif resp.path.type == pb.common_pb2.NOTIFICATION:
                    notif_meta = fc.meta.get_notification(parent.path.meta, resp.path.metaIdent)
                    path = fc.meta.Path(parent.path, notif_meta)
                elif resp.path.type == pb.common_pb2.RPC_INPUT:
                    input_meta = parent.path.meta.input
                    path = fc.meta.Path(parent.path, input_meta)
                    print(f'input_meta, sel_hnd={hnd_id}, num ddefs = {len(path.meta.definitions)} {input_meta}')
                elif resp.path.type == pb.common_pb2.RPC_OUTPUT:
                    output_meta = parent.path.meta.output
                    path = fc.meta.Path(parent.path, output_meta)
                elif resp.path.type == pb.common_pb2.LIST_ITEM:
                    key = None
                    if resp.path.key != None and len(resp.path.key) > 0:
                        key = [fc.val.proto_decode(v) for v in resp.path.key]
                    meta = parent.path.meta
                    path = fc.meta.Path(parent.path, meta, key=key)
                else:
                    raise Exception(f"unrecognized path segment type {resp.path.type} at {parent.path.meta.ident}")
                sel = Selection(driver, hnd_id, node, path, parent, parent.browser, inside_list=resp.insideList) 
            else:
                browser = fc.node.Browser.resolve(driver, resp.browserHnd)
                path = fc.meta.Path(None, browser.module)
                sel = Selection(driver, hnd_id, node, path, None, browser)
        return sel

    def action(self, inputNode=0):
        inputNodeHnd = 0        
        if inputNode != 0:
            inputNodeHnd = resolve_node(self.driver, inputNode).hnd.id
        req = pb.fc_pb2.ActionRequest(selHnd=self.hnd.id, inputNodeHnd=inputNodeHnd)
        resp = self.driver.g_nodes.Action(req)
        outputSel = None
        if resp.outputSelHnd:
            outputSel = Selection.resolve(self.driver, resp.outputSelHnd)
        return outputSel
    
    def notification(self, callback):
        req = pb.fc_pb2.NotificationRequest(selHnd=self.hnd.id)
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

    def upsert_into(self, n):
        n = resolve_node(self.driver, n)
        req = pb.fc_pb2.SelectionEditRequest(op=pb.fc_pb2.UPSERT_INTO, selHnd=self.hnd.id, nodeHnd=n.hnd.id)        
        self.driver.g_nodes.SelectionEdit(req)

    def upsert_from(self, n):
        n = resolve_node(self.driver, n)
        req = pb.fc_pb2.SelectionEditRequest(op=pb.fc_pb2.UPSERT_FROM, selHnd=self.hnd.id, nodeHnd=n.hnd.id)        
        self.driver.g_nodes.SelectionEdit(req)

    def insert_from(self, n):
        n = resolve_node(self.driver, n)
        req = pb.fc_pb2.SelectionEditRequest(op=pb.fc_pb2.INSERT_FROM, selHnd=self.hnd.id, nodeHnd=n.hnd.id)        
        self.driver.g_nodes.SelectionEdit(req)

    def insert_into(self, n):
        n = resolve_node(self.driver, n)
        req = pb.fc_pb2.SelectionEditRequest(op=pb.fc_pb2.INSERT_INTO, selHnd=self.hnd.id, nodeHnd=n.hnd.id)        
        self.driver.g_nodes.SelectionEdit(req)

    def upsert_into_set_defaults(self, n):
        n = resolve_node(self.driver, n)
        req = pb.fc_pb2.SelectionEditRequest(op=pb.fc_pb2.UPSERT_INTO_SET_DEFAULTS, selHnd=self.hnd.id, nodeHnd=n.hnd.id)        
        self.driver.g_nodes.SelectionEdit(req)

    def upsert_from_set_defaults(self, n):
        n = resolve_node(self.driver, n)
        req = pb.fc_pb2.SelectionEditRequest(op=pb.fc_pb2.UPSERT_FROM_SET_DEFAULTS, selHnd=self.hnd.id, nodeHnd=n.hnd.id)        
        self.driver.g_nodes.SelectionEdit(req)

    def update_into(self, n):
        n = resolve_node(self.driver, n)
        req = pb.fc_pb2.SelectionEditRequest(op=pb.fc_pb2.UPDATE_INTO, selHnd=self.hnd.id, nodeHnd=n.hnd.id)        
        self.driver.g_nodes.SelectionEdit(req)

    def update_from(self, n):
        n = resolve_node(self.driver, n)
        req = pb.fc_pb2.SelectionEditRequest(op=pb.fc_pb2.UPDATE_FROM, selHnd=self.hnd.id, nodeHnd=n.hnd.id)        
        self.driver.g_nodes.SelectionEdit(req)

    def replace_from(self, n):
        n = resolve_node(self.driver, n)
        req = pb.fc_pb2.SelectionEditRequest(op=pb.fc_pb2.REPLACE_FROM, selHnd=self.hnd.id, nodeHnd=n.hnd.id)        
        self.driver.g_nodes.SelectionEdit(req)

    def find(self, path):
        req = pb.fc_pb2.FindRequest(selHnd=self.hnd.id, path=path)
        resp = self.driver.g_nodes.Find(req)
        return Selection.resolve(self.driver, resp.selHnd)


def resolve_node(driver, n):
    if not n:
        # nil node
        return 0
    if isinstance(n, int):
        # cached node
        try:
            return fc.handles.Handle.require(driver, n)
        except KeyError:
            return fc.handles.RemoteRef(driver, n)
    if not n.hnd:
        # unregistered local node about to be registered with go
        resp = driver.g_nodes.NewNode(pb.fc_pb2.NewNodeRequest())
        n.hnd = fc.handles.Handle(driver, resp.nodeHnd, n)
    return n


class Browser():

    def __init__(self, driver, module, node=None, node_src=None, hnd_id=None):
        self.driver = driver
        if not hnd_id:
            req = pb.fc_pb2.NewBrowserRequest(moduleHnd=module.hnd.id)
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
        g_req = pb.fc_pb2.BrowserRootRequest(browserHnd=self.hnd.id)
        resp = self.driver.g_nodes.BrowserRoot(g_req)
        return Selection.resolve(self.driver, resp.selHnd)

    @classmethod
    def resolve(cls, driver, hnd_id):
        try:
            return fc.handles.Handle.require(driver, hnd_id)
        except KeyError:
            req = pb.fc_pb2.GetBrowserRequest(browserHnd=hnd_id)
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

class ListRequest():

    def __init__(self, sel, meta, new, delete, row, first, key):
        self.sel = sel
        self.meta = meta
        self.new = new
        self.delete = delete
        self.row = row
        self.first = first
        self.key = key

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


class NodeRequest():

    def __init__(self, sel, new=False, delete=False):
        self.sel = sel
        self.new = new
        self.delete = delete


class XNodeServicer(pb.fc_x_pb2_grpc.XNodeServicer):
    """Bridge between python node navigation and go node navigation"""

    def __init__(self, driver):
        self.driver = driver

    def XNext(self, g_req, context):
        try:
            sel = Selection.resolve(self.driver, g_req.selHnd)
            meta = sel.path.meta
            key_in = None
            if g_req.key != None:
                key_in = []
                for g_key_val in g_req.key:
                    key_in.append(fc.val.proto_decode(g_key_val))

            req = ListRequest(sel, meta, g_req.new, g_req.delete, g_req.row, g_req.first, key_in)
            child, key_out = sel.node.next(req)
            if child != None:
                child = resolve_node(self.driver, child)
                g_key_out = None
                if key_out != None:
                    for v in key_out:
                        g_key_out.append(fc.val.proto_encode(v))
                return pb.fc_x_pb2.XNextResponse(nodeHnd=child.hnd.id, key=g_key_out)
            else:
                return pb.fc_x_pb2.XNextResponse()
        except Exception as error:
            logging.debug(f'XNext error')
            print(traceback.format_exc())
            raise error
        

    def XChild(self, g_req, context):
        try:
            sel = Selection.resolve(self.driver, g_req.selHnd)
            meta = fc.meta.get_def(sel.path.meta, g_req.metaIdent)
            req = ChildRequest(sel, meta, g_req.new, g_req.delete)
            child = sel.node.child(req)
            if child != None:
                child = resolve_node(self.driver, child)
                return pb.fc_x_pb2.XChildResponse(nodeHnd=child.hnd.id)
            else:
                return pb.fc_x_pb2.XChildResponse()
        except Exception as error:
            print(traceback.format_exc())
            raise error

    def XField(self, g_req, context):
        try:
            sel = Selection.resolve(self.driver, g_req.selHnd)
            print(f'xfield, sel.hnd={g_req.selHnd}, node_hnd={sel.node.hnd.id}')
            meta = fc.meta.get_def(sel.path.meta, g_req.metaIdent)
            req = FieldRequest(sel, meta, g_req.write, g_req.clear)
            write_val = None
            if g_req.write:
                write_val = fc.val.proto_decode(g_req.toWrite)
            read_val = sel.node.field(req, write_val)
            if not g_req.write:
                fromRead = fc.val.proto_encode(read_val)
                resp = pb.fc_x_pb2.XFieldResponse(fromRead=fromRead)
            else:
                resp = pb.fc_x_pb2.XFieldResponse()
            return resp
        except Exception as error:
            print(traceback.format_exc())
            raise error

    def XSelect(self, g_req, context):
        # TODO
        pass

    def XAction(self, g_req, context):
        sel = Selection.resolve(self.driver, g_req.selHnd)
        # TODO: Id this inconsistent?
        # meta = fc.meta.require_def(sel.path.meta, g_req.metaIdent)
        meta = sel.path.meta
        input = None
        if g_req.inputSelHnd != 0:
            print(f'xaction, inputSelHnd={g_req.inputSelHnd}')
            input = Selection.resolve(self.driver, g_req.inputSelHnd)
        print(f'about to call action on node {sel.node}')
        output = sel.node.action(ActionRequest(sel, meta, input))
        print(f'done calling action on node {sel.node}')
        outputNodeHnd = None
        if output != None:
            outputNodeHnd = resolve_node(self.driver, output).hnd.id
        return pb.fc_x_pb2.XActionResponse(outputNodeHnd=outputNodeHnd)


    def XNodeSource(self, g_req, context):
        browser = fc.node.Browser.resolve(self.driver, g_req.browserHnd)
        n = resolve_node(self.driver, browser.node())
        return pb.fc_x_pb2.XNodeSourceResponse(nodeHnd=n.hnd.id)


    def XNotification(self, g_req, context):
        q = queue.Queue()
        sel = Selection.resolve(self.driver, g_req.selHnd)
        meta = sel.path.meta
        closer = sel.node.notification(NotificationRequest(sel, meta, q))
        def stream_closed():
            q.put(None)
        context.add_callback(stream_closed)
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
    
    def XBeginEdit(self, g_req, context):
        sel = Selection.resolve(self.driver, g_req.selHnd)
        sel.node.begin_edit(NodeRequest(sel, new=g_req.new, delete=g_req.delete))
        return pb.fc_x_pb2.XBeginEditResponse()

    def XEndEdit(self, g_req, context):
        sel = Selection.resolve(self.driver, g_req.selHnd)
        sel.node.end_edit(NodeRequest(sel, new=g_req.new, delete=g_req.delete))
        return pb.fc_x_pb2.XEndEditResponse()

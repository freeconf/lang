from fc.nodeutil import reflect, extend
import time
import threading


# Simple application, no connection to management 
class Car():

    def __init__(self):
        self.speed = 0
        self.miles = 0
        self.running = False
        self.thread = None
        self.listeners = []

    def start(self, running):
        if self.running != running:
            self.running = running
            if running:
                self.thread = threading.Thread(target=self.run, name="Car")
                self.thread.start()
            if not running:
                self.thread = None

    def reset(self):
        self.miles = 0

    def run(self):
        self.update_listeners("started")
        while self.running:
            time.sleep(0.01)
            self.miles = self.miles + self.speed
        self.update_listeners("stopped")

    def on_update(self, listener):
        self.listeners.append(listener)
        def closer():
            self.listeners.remove(listener)
        return closer

    def update_listeners(self, event):
        for l in self.listeners:
            l(event)

# Bridge car to FC management library
def manage(c):

    def action(node, req):
        if req.meta.ident == 'stop':
            c.start(False)
        elif req.meta.ident == 'start':
            c.start(True)
        else:
            return node.action(req)
        return None

    def notification(node, req):
        if req.meta.ident == 'update':
            def listener(event):
                req.send(reflect.Reflect({
                    "event": event
                }))
            closer = c.on_update(listener)
            return closer
        
        return node.notification(req)

    # because car's members and methods align with yang, we can use 
    # reflection for all of the CRUD
    return extend.Extend(
        base = reflect.Reflect(c),
        on_action = action, on_notification=notification)

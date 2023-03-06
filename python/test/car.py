import fc.nodeutil
import time
import threading

# Simple application, no connection to management 
class Car():

    def __init__(self):
        self.speed = 0
        self.miles = 0
        self.running = False
        self.thread = None

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
        while self.running:
            time.sleep(0.01)
            self.miles = self.miles + self.speed

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

    # because car's members and methods align with yang, we can use 
    # reflection for all of the CRUD
    return fc.nodeutil.Extend(
        base = fc.nodeutil.Reflect(c),
        on_action = action)
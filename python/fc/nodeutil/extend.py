class Extend():

    def __init__(self, base, on_child=None, on_field=None, on_action=None, on_notification=None):
        self.hnd = 0
        self.base = base
        self.on_child = on_child
        self.on_field = on_field
        self.on_action = on_action
        self.on_notification = on_notification


    def child(self, r):
        if self.on_child:
            return self.on_child(self, self.base, r)
        return self.base.child(r)


    def field(self, r, write_val):
        if self.on_field:
            return self.on_field(self, self.base, r, write_val)
        return self.base.field(r, write_val)


    def action(self, r):
        if self.on_action:
            return self.on_action(self.base, r)
        return self.base.action(r)


    def notification(self, r):
        if self.on_notification:
            return self.on_notification(self.base, r)
        return self.base.notification(r)
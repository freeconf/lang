import pb.fc_g_pb2


# a pointer to an object in Go that is being held for python
# until this handle is no longer referenced.  
class Handle():

    def __init__(self, driver, id, obj):
        self.id = id
        self.driver = driver
        driver.handles[id] = obj

    def release(self):
        # if driver is unloaded, no need to release anything as it would be killed
        # with the fc-lang process
        if self.driver.handles:
            self.driver.g_handles.Release(pb.fc_g_pb2.ReleaseRequest(hnd=self.id))

    # TODO Get's claimed to early
    # def __del__(self):
    #     self.release()

    @classmethod
    def lookup(cls, driver, id):
        return driver.handles.get(id, None)

    @classmethod
    def require(cls, driver, id):
        try:
            return driver.handles[id]
        except KeyError:
            raise KeyError(f'could not resolve hnd {id}')

# For handles that never are used in python except to be passed back to
# go.  
class RemoteRef():

    def __init__(self, driver, id):
        self.hnd = Handle(driver, id, self)


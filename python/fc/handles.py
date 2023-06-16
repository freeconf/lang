import fc.pb.fc_pb2

# For handles that never are used in python except to be passed back to
# go.  
class RemoteRef():

    def __init__(self, driver, id):
        self.hnd = id
        driver.obj_strong.store_hnd(id, self)


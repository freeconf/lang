
handles = {}
objects = {}
# Go get's even handles, X get's odd so they never collide
counter = 1

def get(hnd):
    global objects    
    return objects[hnd]

def put(obj, hnd=None):
    global handles
    global objects
    global counter
    if not hnd:
        counter, hnd = counter + 2, counter
    objects[hnd] = obj
    handles[obj] = hnd
    return hnd

def hnd(obj):
    return handles[obj]

def release_hnd(hnd):
    global handles
    obj = objects.pop(hnd)
    if obj:
        del handles[obj]

def release_obj(obj):
    global handles
    hnd = handles.pop(obj)
    if hnd:
        del objects[hnd]

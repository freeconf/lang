
handles = {}
counter = 1

def get(id):
    global handles
    return handles[id]

def put(obj):
    global handles
    global counter
    counter, id = counter + 1, counter
    handles[id] = obj
    return id

def release(id):
    global handles
    del handles[id]

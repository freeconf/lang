package org.freeconf.fc.meta;

import org.freeconf.fc.driver.Closeable;
import org.freeconf.fc.driver.Driver;

public class Module implements Closeable {
    private long poolId;
    private String ident;
    private String desc;

    public Module(long myPoolId, String myident, String mydesc) {
        poolId = myPoolId;
        ident = myident;
        desc = mydesc;
    }

    public String GetIdent()  {
        return ident;
    }

    public String GetDesc() {
        return desc;
    }

    public void Close() {
        org.freeconf.fc.driver.Driver.release(poolId);
    }

    static {
        Driver.loadLibrary();
    }
}

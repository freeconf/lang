package org.freeconf.fc;

public class Module {
    private String ident;
    private String desc;

    public Module(String myident, String mydesc) {
        ident = myident;
        desc = mydesc;
    }

    public String GetIdent()  {
        return ident;
    }

    public String GetDesc() {
        return desc;
    }
}

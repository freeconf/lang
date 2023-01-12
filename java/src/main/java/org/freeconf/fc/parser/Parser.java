package org.freeconf.fc.parser;

import org.freeconf.fc.meta.Module;
import org.freeconf.fc.driver.Driver;

public class Parser {
    public native Module parse(String ypath, String yfile);

    static {
        Driver.loadLibrary();
    }
}
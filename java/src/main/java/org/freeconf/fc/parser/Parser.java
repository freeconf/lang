package org.freeconf.fc.parser;

import org.freeconf.fc.meta.Module;

public class Parser {

    public native Module parse(String ypath, String yfile);
}
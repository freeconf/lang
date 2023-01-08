package org.freeconf.fc;

import static org.junit.Assert.assertEquals;

import org.junit.Test;

public class ParserTest {

    @Test
    public void shouldParse() {
        String ypath = System.getenv("YANGPATH");
        Module m = new Parser().parseFile(ypath, "testme");
        assertEquals("testme", m.GetIdent());
    }    
}

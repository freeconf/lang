package org.freeconf.fc;

import static org.junit.Assert.assertEquals;

import org.junit.Test;
import org.freeconf.fc.parser.Parser;
import org.freeconf.fc.driver.Driver;
import org.freeconf.fc.meta.Module;

public class ParserTest {

    @Test
    public void shouldParse() {
        new Driver();
        String ypath = System.getenv("YANGPATH");
        Module m = new Parser().parse(ypath, "testme");
        assertEquals("testme", m.GetIdent());
        m.Close();
    }    
}

package org.freeconf.fc.driver;

/**
 * Object that require closing to release underlying driver
 */
public interface Closeable {
    void Close();
}

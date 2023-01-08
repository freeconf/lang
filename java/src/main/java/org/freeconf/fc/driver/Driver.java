package org.freeconf.fc.driver;

public class Driver {
    static {
        try {
            // Exampple:
            //  export LD_LIBRARY_PATH=/home/joe/lang/c/lib:/home/joe/lang/java/target/so
            System.loadLibrary("fc-j");
        } catch (UnsatisfiedLinkError e) {
            System.err.println("Failed to load freeconf library.\n" + e);
            System.exit(1);
        }
    }

    public Driver() {        
    }
    
    public static native void release(long poolId);
}
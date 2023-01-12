package org.freeconf.fc.driver;

/**
 * Low-level coordination with underlying JNI connection to FreeCONF C library
 * which in turn has a connection to Go library
 */
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

    /**
     * NOOP method but adding this to each class that used the JNI layer will ensure the
     * static block above is exercised.
     */
    public static void loadLibrary() {}

    /**
     * Release underlying memory. Do not call twice for same poolId.  Internal
     * only FreeCONF Java library should call this
     * 
     * @param poolId - Should be given to you by underlying call
     */
    public static native void release(long poolId);
}
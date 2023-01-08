package org.freeconf.fc;

public class Parser {
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

    private native Module parse(String ypath, String yfile);

    public Module parseFile(String ypath, String yfile) {
        return parse(ypath, yfile);
    }
}
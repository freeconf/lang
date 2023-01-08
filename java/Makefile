FC_C_DIR = $(abspath ../c)
DEST_DIR = target/so

all: build link

build:
	test -d $(DEST_DIR) || mkdir -p $(DEST_DIR)
	javac -h src/main/c src/main/java/org/freeconf/fc/Parser.java
	gcc -c -fPIC \
		-I$(JAVA_HOME)/include -I${JAVA_HOME}/include/linux \
		-I$(FC_C_DIR)/lib \
		src/main/c/parser.c \
		-o $(DEST_DIR)/parser.o

link:
	gcc -shared -fPIC \
		-L$(FC_C_DIR)/lib \
		-o $(DEST_DIR)/libfc-j.so \
		$(DEST_DIR)/parser.o \
		-lc -lfc

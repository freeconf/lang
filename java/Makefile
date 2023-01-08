FC_C_DIR = $(abspath ../c)
DEST_DIR = ./target/so

LIB = libfc-j.so

all: build

JNI_SRC_JAVA = \
	org/freeconf/fc/parser/Parser.java \
	org/freeconf/fc/driver/Driver.java

JNI_SRC_H = \
	$(foreach F,$(JNI_SRC_JAVA),$(basename $(subst /,_,$F)).h)

JNI_SRC_C = \
	$(JNI_SRC_H:.h=.c)

OBJS = \
	$(JNI_SRC_C:.c=.o)

build: $(DEST_DIR) $(DEST_DIR)/$(LIB)

$(DEST_DIR):
	mkdir -p $(DEST_DIR)

$(DEST_DIR)/%.h : $(foreach F,$(JNI_SRC_JAVA),src/main/java/$(F))
	javac -h $(DEST_DIR) --source-path src/main/java $^

# $(DEST_DIR)/%.h : 
# 	echo here

src/main/c/%.c : $(DEST_DIR)/%.h

$(DEST_DIR)/%.o : src/main/c/%.c
	gcc -c -fPIC \
		-I$(JAVA_HOME)/include -I${JAVA_HOME}/include/linux \
		-I$(DEST_DIR) -I$(FC_C_DIR)/lib \
		$< \
		-o $@

$(DEST_DIR)/$(LIB): $(foreach F,$(OBJS),$(DEST_DIR)/$(F))
	gcc -shared -fPIC \
		-L$(DEST_DIR) -L$(FC_C_DIR)/lib \
		-o $@ $^ \
		-lc -lfc

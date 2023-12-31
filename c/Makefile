FC_LANG_DIR = ..
SRC = $(wildcard ./*.c)
OBJS = $(foreach F,$(SRC),./out/$(notdir $(F:.c=.o)))

TESTS = \
	out/test/test_parser \
	out/test/test_node

export LD_LIBRARY_PATH = $(FC_LANG_DIR)/lib:./lib
export YANGPATH=$(abspath ../yang)

# I think cgo generates code that triggers 'stack smashing detected'
# so i had to disable this check
GCC_FLAGS = \
	-fno-stack-protector \
	-fPIC

LIB_DIRS = \
	-L$(FC_LANG_DIR)/lib \
	-L/usr/local/x86_64-linux-gnu

LIBS = \
	-lfc \
	-lcbor

INCLUDE_DIRS = \
	-I$(FC_LANG_DIR)/lib \
	-I./include \
	-I.

all : lib/libfc-c.so test

.PHONY: test
test: $(TESTS)

lib/libfc-c.so: $(OBJS)
	gcc \
	$(LIB_DIRS) \
	-shared -o $@ $^ \
	$(LIBS)

out/test/% : test/%.c
	! test -d out/test || mkdir -p out/test
	gcc \
		$(GCC_FLAGS) \
		-Wall \
		$(INCLUDE_DIRS) \
		$(LIB_DIRS) -L./lib \
		-o $@ $^ \
		$(LIBS) -lfc-c
	$@

out/%.o : %.c
	gcc -c \
		$(GCC_FLAGS) \
		-Wall \
		$(INCLUDE_DIRS) \
		$^ -o $@

GO_SRC = $(wildcard ./go/*.go ./go/*/*.go)
TEST_SRC = $(wildcard ./test/*.c)
TESTS = \
	test_parser
TEST_OBJ = $(foreach T,$(TESTS), out/$(T).out)
export LD_LIBRARY_PATH = ./lib
export YANGPATH = $(abspath ../yang)

# I think cgo generates code that triggers 'stack smashing detected'
# so i had to disable this check
GCC_FLAGS = -fno-stack-protector 

LIB_DIRS = \
	-L./lib \
	-L/usr/local/x86_64-linux-gnu \
	-L/home/dhubler/vendor/QCBOR/build

LIBS = \
	-lfc \
	-lcbor \
	-l:libqcbor.a -lm 

INCLUDE_DIRS = \
	-I./lib \
	-I/home/dhubler/vendor/QCBOR/inc

all : lib/libfc.so test

test: $(TESTS)

$(TESTS) : % : out/%.out
	./out/$@.out

lib/libfc.so: $(GO_SRC)
		go build \
			-buildmode=c-shared \
			-o lib/libfc.so ./go/*.go

out/%.out : $(TEST_SRC)
	gcc \
		$(GCC_FLAGS) \
		-Wall \
		$(INCLUDE_DIRS) \
		$(LIB_DIRS) \
		./test/$*.c -o $@ \
		$(LIBS)

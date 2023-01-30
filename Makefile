export LD_LIBRARY_PATH = out
export YANGPATH=$(abspath test/yang)

# I think cgo generates code that triggers 'stack smashing detected'
# so i had to disable this check
GCC_FLAGS = \
	-fno-stack-protector \
	-fsanitize=address \
	-fPIC

INCLUDE_DIRS = \
	-I. -I./out

LIB_DIRS = \
	-L./out \
	-L/usr/local/x86_64-linux-gnu

LIBS = \
	-lfc \
	-lcbor

TESTS = \
	out/test_parser \
	out/test_node

all : generate lib test

generate:
	cd codegen; \
		go generate .

.PHONY: lib
lib : out/libfc.so

.PHONY: out/libfc.so
out/libfc.so:
	test -d out || mkdir out
	CGO_ENABLED=1 \
		go build \
		-buildmode=c-shared \
		-o out/libfc.so .

.PHONY: test
test: $(TESTS)

out/test_% : test/test_%.c generate out/libfc.so 
	gcc \
		$(GCC_FLAGS) \
		-Wall \
		$(INCLUDE_DIRS) \
		$(LIB_DIRS) \
		-o $@ $< \
		$(LIBS)
	$@
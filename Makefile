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
	go run codegen/code_gen_main.go -codegen_dir ./codegen \
		./*.in \
		python/fc/*.in

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

out/test_% : test/test_%.c
	gcc \
		$(GCC_FLAGS) \
		-Wall \
		$(INCLUDE_DIRS) \
		$(LIB_DIRS) \
		-o $@ $< \
		$(LIBS)
	$@

proto:
	test -d comm/pb || mkdir comm/pb
	protoc \
		--go_out=comm  \
		--go-grpc_out=comm \
		comm/*.proto

proto-py:
	test -d python/pb || mkdir python/pb
	cd python; \
		. venv/bin/activate && \
		python -m grpc_tools.protoc \
			-I../comm --python_out=pb --pyi_out=pb \
			--grpc_python_out=pb ../comm/meta.proto
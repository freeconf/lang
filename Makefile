export YANGPATH = $(abspath test/yang)
export PATH := $(PATH):$(abspath ./bin)
export PYTHONPATH := $(abspath python)

all : generate proto bin test

generate:
	go run codegen/main.go \
		./*.in \
		python/freeconf/*.in

.PHONY: test
test: test-go test-py

.PHONY: dist
proto: proto-go proto-py

#################
## G O
#################
test-go:
	FC_LANG=GO go test . ./...

deps-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

.PHONY: bin
bin : bin/fc-lang bin/fc-lang-dbg

.PHONY: bin/fc-lang
# see https://www.jetbrains.com/help/go/attach-to-running-go-processes-with-debugger.html#step-2-build-the-application
bin/fc-lang-dbg : BUILD_OPTS=-gcflags="all=-N -l"
bin/fc-lang-dbg bin/fc-lang :
	test -d $(dir $@) || mkdir -p $(dir $@)
	go build $(BUILD_OPTS) -o $@ cmd/fc-lang/main.go


proto-go:
	! test -d pb || rm -rf pb
	mkdir pb
	protoc \
		-I./proto \
		--go_out=. \
		--go-grpc_out=. \
		./proto/freeconf/pb/*.proto

#################
## P Y T H O N
#################

proto-py:
	! test -d python/freeconf/pb || rm -rf python/freeconf/pb
	mkdir python/freeconf/pb
	cd python; \
		python3 -m grpc_tools.protoc \
			-I../proto \
			--python_out=. \
			--pyi_out=. \
			--grpc_python_out=. \
			../proto/freeconf/pb/*.proto

# Loosely ordered by lower level to to higher level operations
PY_TESTS = \
	test_val.py \
	test_driver.py \
	test_node.py \
	test_reflect.py \
	test_parser.py \
	test_car.py \
	test_restconf.py

test-py:
	cd python/tests; \
		$(foreach T,$(PY_TESTS),echo $(T) && python3 $(T) || exit;)
	FC_LANG=GO go test ./test

deps-py:
	pip install build
	cd python; \
		pip install -e . && \
		pip install -e ".[dev]"

dist-py :
	cd python; \
		python3 -m build

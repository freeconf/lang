VER = 0.1.0
export YANGPATH = $(abspath test/yang)
export PATH := $(PATH):$(abspath ./bin)
export PYTHONPATH := $(abspath python)

all : generate proto dist-go test dist-py

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
	FC_LANG=go go test . ./...

deps-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2


# Just the popular ones.  You can easily build binary for missing platform
PLATFORMS = \
  darwin-amd64 \
  darwin-arm64 \
  linux-amd64 \
  windows-amd64

BIN_TARGETS = $(foreach P,$(PLATFORMS), bin/fc-lang-$(VER)-$(P))

dist-go: $(BIN_TARGETS)

bin/fc-lang-$(VER)-darwin-amd64: BUILD_ENV=GOARCH=amd64 GOOS=darwin
bin/fc-lang-$(VER)-darwin-arm64: BUILD_ENV=GOARCH=amd64 GOOS=darwin
bin/fc-lang-$(VER)-windows-amd64: BUILD_ENV=GOARCH=amd64 GOOS=windows
bin/fc-lang-$(VER)-windows-amd64: BIN_EXT=.exe
bin/fc-lang-$(VER)-linux-amd64: BUILD_ENV=GOARCH=amd64 GOOS=linux

# see https://www.jetbrains.com/help/go/attach-to-running-go-processes-with-debugger.html#step-2-build-the-application
bin/fc-lang-dbg : BUILD_OPTS=-gcflags="all=-N -l"

debug:
	echo $(BIN_TARGETS)

.PHONY: bin/fc-lang bin/fc-lang-dbg $(BIN_TARGETS)

bin/fc-lang bin/fc-lang-dbg $(BIN_TARGETS):
	test -d $(dir $@) || mkdir -p $(dir $@)
	$(BUILD_ENV) go build $(BUILD_OPTS) -o $@$(BIN_EXT) cmd/fc-lang/main.go

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
	touch python/freeconf/pb/__init__.py
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
	test_parser.py \
	test_car.py \
	test_restconf.py \
	test_util_node.py \
	test_node_action.py

test-py:
	cd python/tests; \
		$(foreach T,$(PY_TESTS),echo $(T) && python3 $(T) || exit;)
	FC_LANG=python go test ./test

deps-py:
	pip install build
	cd python; \
		pip install -e . && \
		pip install -e ".[dev]"

dist-py :
	! test -d python/freeconf.egg-info || rm -rf python/freeconf.egg-info
	cd python; \
		python3 -m build 

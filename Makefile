export YANGPATH = $(abspath test/yang)
export PATH := $(PATH):$(abspath ./bin)
export PYTHONPATH := $(abspath python)

all : generate proto bin test test-py

generate:
	go run codegen/main.go \
		./*.in \
		python/freeconf/*.in

.PHONY: test
test: test-go test-py

test-go:
	go test . ./...

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
		. ../venv/bin/activate && \
		$(foreach T,$(PY_TESTS),echo $(T) && python $(T) || exit;)

.PHONY: bin
bin : bin/fc-lang bin/fc-lang-dbg

.PHONY: bin/fc-lang
# see https://www.jetbrains.com/help/go/attach-to-running-go-processes-with-debugger.html#step-2-build-the-application
bin/fc-lang-dbg : BUILD_OPTS=-gcflags="all=-N -l"
bin/fc-lang-dbg bin/fc-lang :
	test -d $(dir $@) || mkdir -p $(dir $@)
	go build $(BUILD_OPTS) -o $@ cmd/fc-lang/main.go

proto: proto-go proto-py

install-deps:
	sudo mkdir /opt/protoc
	sudo unzip -d /opt/protoc/ ./protoc-21.12-linux-x86_64.zip
	export PATH=$PATH:/opt/protoc/bin
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2	

proto-go:
	! test -d pb || rm -rf pb
	mkdir pb
	protoc \
		-I./proto \
		--plugin=protoc-gen-go-grpc=${HOME}/go/bin/protoc-gen-go-grpc \
		--plugin=protoc-gen-go=${HOME}/go/bin/protoc-gen-go \
		--go_out=. \
		--go-grpc_out=. \
		./proto/freeconf/pb/*.proto

proto-py:
	! test -d python/freeconf/pb || rm -rf python/freeconf/pb
	mkdir python/freeconf/pb
	cd python; \
		. venv/bin/activate && \
		python -m grpc_tools.protoc \
			-I../proto \
			--python_out=. \
			--pyi_out=. \
			--grpc_python_out=. \
			../proto/freeconf/pb/*.proto

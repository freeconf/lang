export YANGPATH = $(abspath test/yang)
export PATH := $(PATH):./bin
export PYTHONPATH := $(abspath python):$(abspath python)/pb

all : generate bin proto test test-py

generate:
	go run codegen/main.go \
		./*.in \
		python/fc/*.in

.PHONY: test
test: test-go test-py

test-go:
	go test .

PY_TESTS = \
	test_val.py \
	test_driver.py \
	test_node.py \
	test_car.py \
	test_parser.py \
	test_restconf.py

test-py:
	cd python; \
		. venv/bin/activate && \
		$(foreach T,$(PY_TESTS),python test/$(T) || exit;)


.PHONY: bin
bin : bin/fc-lang bin/fc-lang-dbg

.PHONY: bin/fc-lang
# see https://www.jetbrains.com/help/go/attach-to-running-go-processes-with-debugger.html#step-2-build-the-application
bin/fc-lang-dbg : BUILD_OPTS=-gcflags="all=-N -l"
bin/fc-lang-dbg bin/fc-lang :
	test -d $(dir $@) || mkdir -p $(dir $@)
	go build $(BUILD_OPTS) -o $@ cmd/fc-lang/main.go


proto: proto-go proto-py

proto-go:
	! test -d pb || rm -rf pb
	mkdir pb
	protoc \
		-I./proto \
		--go_out=. \
		--go-grpc_out=. \
		./proto/*.proto

proto-py:
	! test -d python/pb || rm -rf python/pb
	mkdir python/pb
	cd python; \
		. venv/bin/activate && \
		python -m grpc_tools.protoc \
			-I../proto \
			--python_out=pb \
			--pyi_out=pb \
			--grpc_python_out=pb \
			../proto/*.proto

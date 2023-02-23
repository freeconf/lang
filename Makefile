export YANGPATH = $(abspath test/yang)
export PATH := $(PATH):./bin

all : generate lib test

generate:
	go run codegen/main.go \
		./*.in \
		python/fc/*.in

.PHONY: test
test:
	go test .

PY_TESTS = \
	test_driver.py \
	test_parser.py \
	test_node.py

test-py:
	cd python; \
		. venv/bin/activate && \
		python ./test_

.PHONY: bin/fc-lang
bin/fc-lang :
	test -d $(dir $@) || mkdir -p $(dir $@)
	go build -o $@ cmd/fc-lang/main.go


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

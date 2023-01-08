export CGO_ENABLED = 1

GRPC_DEST = fclang_grpc

clib:
		go build \
			-buildmode=c-shared \
			-o c/lib/libfc.so c/go/*.go

proto:
	protoc \
		--go_out=$(GRPC_DEST) \
		--go_opt=paths=source_relative \
    	--go-grpc_out=$(GRPC_DEST) \
		--go-grpc_opt=paths=source_relative \
    	fc-lang.proto

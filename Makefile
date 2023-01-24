export CGO_ENABLED = 1

all : generate lib

generate:
	cd emeta; \
		go generate .

.PHONY: lib
lib : lib/libfc.so

.PHONY: lib/libfc.so
lib/libfc.so:
		go build \
			-buildmode=c-shared \
			-o lib/libfc.so *.go

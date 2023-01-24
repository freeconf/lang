export CGO_ENABLED = 1

all : generate lib

SRC = \
	path.go \
	browser.go \
	err.go \
	main.go \
	meta.go \
	node.go \
	nodeutil.go \
	parser.go \
	pool.go \
	selection.go \
	val.go


generate:
	cd emeta; \
		go generate .

.PHONY: lib
lib : lib/libfc.so

.PHONY: lib/libfc.so
lib/libfc.so:
		go build \
			-buildmode=c-shared \
			-o lib/libfc.so $(SRC)

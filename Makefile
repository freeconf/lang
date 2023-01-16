export CGO_ENABLED = 1

lib/libfc.so:
		go build \
			-buildmode=c-shared \
			-o lib/libfc.so *.go

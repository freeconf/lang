## Compiler

```
sudo apt install gcc g++
```

## CBOR

```
sudo apt install libcbor-dev
```

## Flatbuffers
https://google.github.io/flatbuffers/flatbuffers_guide_building.html

```
sudo apt install cmake
cmake -G "Unix Makefiles" -DCMAKE_BUILD_TYPE=Release
make
```

## Setting up protoc

Download protoc from https://github.com/protocolbuffers/protobuf/releases

```
sudo mkdir /opt/protoc
sudo unzip -d /opt/protoc/ ./protoc-21.12-linux-x86_64.zip
export PATH=$PATH:/opt/protoc/bin
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

### Zig

Install `xz-utils` package
Download from https://ziglang.org/download/ and install into /usr/local and put `zg` in `PATH`

## Debugging

```bash
# Which binary to run for language support for FreeCONF's core engine
FC_LANG_EXEC=fc-lang

# Opens a port to listen for Go's Delve debugger on port 999
FC_LANG_DBG_ADDR=:9999
```

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


## C


### Naming conventions

```
fc_meta
fc_meta_path
fc_meta_module
fc_meta_container
fc_meta_list
fc_meta_find
fc_meta_get_ident
fc_meta_get_desc
fc_meta_get_defs
fc_meta_get_exts
fc_meta_array
fc_meta_ext
fc_meta_choice
fc_meta_choice_case
fc_meta_ext_array
fc_meta_ext_def
fc_meta_ext_def_arg
fc_meta_optional_bool

fc_val
fc_val_hnd
fc_val_type

fc_node
fc_node_child
fc_node_field
fc_node_action
fc_node_field_req
fc_node_child_req
fc_node_action_req
fc_node_notify_req
fc_meta_type

fc_error

fc_browser
fc_browser_new

fc_yang_parse

fc_browser_root_select

fc_mem_free

fc_pack_err
fc_pack_decode_fc_meta_module
// other decode methods internal to that do not matter

fc_select
fc_select_upsert_from
fc_select_upsert_to
fc_select_insert_from
fc_select_insert_to
fc_select_delete
fc_json_rdr
fc_json_wtr
```
export GRPC_TRACE=all
GRPC_VERBOSITY=DEBUG

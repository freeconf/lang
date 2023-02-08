package codegen

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/yoheimuta/go-protoparser"
)

func TestParseMetaProto(t *testing.T) {
	reader, err := os.Open("../comm/meta.proto")
	if err != nil {
		t.Fatalf("failed to open, err %v\n", err)
	}
	defer reader.Close()

	got, err := protoparser.Parse(reader)
	if err != nil {
		t.Fatalf("failed to parse, err %v\n", err)
	}

	gotJSON, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal, err %v\n", err)
	}
	fmt.Print(string(gotJSON))
}

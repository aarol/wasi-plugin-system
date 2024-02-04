package main

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"

	"github.com/aarol/wasi-plugin-system/gen/proto"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	protobuf "google.golang.org/protobuf/proto"
)

//go:embed wasm/target/wasm32-wasi/release/plugin.wasm
var plugin []byte

func main() {
	r := wazero.NewRuntime(context.Background())

	defer r.Close(context.Background())

	wasi_snapshot_preview1.MustInstantiate(context.Background(), r)

	req := proto.Request{Input: "Hello, Wasm"}

	b, err := protobuf.Marshal(&req)
	if err != nil {
		panic(err)
	}

	stdin := bytes.NewReader(b)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	config := wazero.NewModuleConfig().WithStdin(stdin).WithStdout(&stdout).WithStderr(&stderr)

	_, err = r.InstantiateWithConfig(context.Background(), plugin, config)
	if err != nil {
		panic(err)
	}
	if stderr.Len() > 0 {
		fmt.Println("Stderr: ", stderr.String())
		return
	}

	var res proto.Response
	err = protobuf.Unmarshal(stdout.Bytes(), &res)
	if err != nil {
		panic(err)
	}
	fmt.Println("Response:", res.Output)
}

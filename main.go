package main

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"

	"github.com/aarol/wasi-plugin-system/gen/plugin"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"google.golang.org/protobuf/encoding/protodelim"
)

//go:embed wasm/target/wasm32-wasi/release/plugin.wasm
var syntaxPlugin []byte

func main() {
	r := wazero.NewRuntime(context.Background())

	defer r.Close(context.Background())

	wasi_snapshot_preview1.MustInstantiate(context.Background(), r)

	stdin := bufio.NewReadWriter()
	var stdout bufio.ReadWriter
	var stderr bufio.ReadWriter
	config := wazero.NewModuleConfig().WithStdin(&stdin).WithStdout(&stdout).WithStderr(&stderr)

	_, err := r.InstantiateWithConfig(context.Background(), syntaxPlugin, config)
	if err != nil {
		panic(err)
	}

	var info plugin.PluginInfo

	protodelim.UnmarshalFrom(bufio.NewReader(&stdout), &info)
	fmt.Println(info.Events)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

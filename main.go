package main

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io"
	"time"

	"github.com/aarol/wasi-plugin-system/gen/plugin"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"google.golang.org/protobuf/proto"
)

//go:embed wasm/target/wasm32-wasi/release/plugin.wasm
var syntaxPlugin []byte

func main() {
	r := wazero.NewRuntime(context.Background())

	defer r.Close(context.Background())

	wasi_snapshot_preview1.MustInstantiate(context.Background(), r)

	req := plugin.Request{
		Req: &plugin.Request_SyntaxRequest{
			SyntaxRequest: &plugin.SyntaxRequest{
				Code: "let a = 56;", Language: "rs",
			},
		},
	}
	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()
	var stderr bytes.Buffer

	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("Panic")
			fmt.Println(stderr.String())
		}
	}()

	config := wazero.NewModuleConfig().WithStdin(stdinReader).WithStdout(stdoutWriter).WithStderr(&stderr)

	compiled, err := r.CompileModule(context.Background(), syntaxPlugin)
	if err != nil {
		panic(err)
	}
	// now := time.Now()
	_, err = r.InstantiateModule(context.Background(), compiled, config)
	if err != nil {
		panic(err)
	}

	b, err := proto.Marshal(&req)
	must(err)

	x := plugin.PluginInfo{
		Events: []plugin.Events{plugin.Events_SYNTAX_HIGHLIGHT},
	}

	b, err = proto.Marshal(&x)
	must(err)
	_, err = stdinWriter.Write(b)
	must(err)

	if stderr.Len() > 0 {
		fmt.Println("Stderr:", stderr.String())
	}
	var res plugin.SyntaxResponse
	must(proto.Unmarshal(stdoutReader, &res))

	fmt.Println("Elapsed", time.Since(now).Milliseconds())
	fmt.Println(res.Output)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

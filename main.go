package main

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/aarol/wasi-plugin-system/gen/plugin"
	"github.com/samber/lo"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

//go:embed rust/target/wasm32-wasi/release/plugin.wasm
var rustWasm []byte

//go:embed assemblyscript/build/release.wasm
var ascWasm []byte

//go:embed go/main.wasm
var goWasm []byte

func main() {
	r := wazero.NewRuntime(context.Background())

	defer r.Close(context.Background())

	wasi_snapshot_preview1.MustInstantiate(context.Background(), r)

	plugins := []*WasmPlugin{}

	goPlugin := lo.Must(NewWasmPlugin(r, goWasm))
	plugins = append(plugins, goPlugin)

	for _, p := range plugins {

		req := plugin.Request{
			Req: &plugin.Request_SyntaxRequest{
				SyntaxRequest: &plugin.SyntaxRequest{
					Code:     "var a = 1",
					Language: "go",
				},
			},
		}
		var res plugin.SyntaxResponse

		lo.Must0(p.Call(&req, &res))

		fmt.Println(res.Output)
	}
}

type WasmPlugin struct {
	module  wazero.CompiledModule
	runtime wazero.Runtime
}

func NewWasmPlugin(runtime wazero.Runtime, wasm []byte) (*WasmPlugin, error) {
	module, err := runtime.CompileModule(context.Background(), wasm)
	if err != nil {
		return nil, err
	}

	return &WasmPlugin{
		runtime: runtime,
		module:  module,
	}, nil
}

func (p *WasmPlugin) Call(req *plugin.Request, res protoreflect.ProtoMessage) (err error) {
	b, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	stdin := bytes.NewReader(b)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	config := wazero.NewModuleConfig().WithStdin(stdin).WithStdout(&stdout).WithStderr(&stderr).WithEnv("HOST_VERSION", "1.0.0")

	defer func() {
		recoveredErr := recover()
		if recoveredErr != nil {
			if stderr.Len() > 0 {
				err = errors.New(stderr.String())
			} else {
				err = fmt.Errorf("%s", recoveredErr)
			}
		}
	}()
	_, err = p.runtime.InstantiateModule(context.Background(), p.module, config)
	if err != nil {
		return err
	}

	if stderr.Len() > 0 {
		return errors.New(stderr.String())
	}

	err = proto.Unmarshal(stdout.Bytes(), res)
	if err != nil {
		return err
	}
	return nil
}

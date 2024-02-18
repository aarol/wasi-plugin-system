package main

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/aarol/wasi-plugin-system/gen/plugin"
	"github.com/samber/lo"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

//go:embed rust/target/wasm32-wasi/release/plugin.wasm
var rustPlugin []byte

//go:embed assemblyscript/build/release.wasm
var ascPlugin []byte

func main() {
	r := wazero.NewRuntime(context.Background())

	defer r.Close(context.Background())

	wasi_snapshot_preview1.MustInstantiate(context.Background(), r)

	wasmPlugin := lo.Must(NewWasmPlugin(r, ascPlugin))

	req := &plugin.Request{Req: &plugin.Request_SyntaxRequest{
		SyntaxRequest: &plugin.SyntaxRequest{
			Code:     "let a = 1;",
			Language: "rs",
		},
	}}

	now := time.Now()
	res := plugin.SyntaxResponse{}
	lo.Must0(wasmPlugin.Call(req, &res))
	fmt.Println(res.Output)
	fmt.Println(time.Since(now).Microseconds())

	req2 := &plugin.Request{Req: &plugin.Request_VersionRequest{
		VersionRequest: &plugin.VersionRequest{},
	}}
	for i := 0; i < 100; i++ {
		now = time.Now()
		var res2 plugin.VersionResponse
		lo.Must0(wasmPlugin.Call(req2, &res2))

		fmt.Println(res2.Version)

		fmt.Println(time.Since(now).Microseconds())
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

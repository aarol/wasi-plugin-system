package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/aarol/wasi-plugin-system/go/gen/plugin"
	"github.com/alecthomas/chroma/v2/quick"
	"google.golang.org/protobuf/proto"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute plugin: %s", err)
	}
}

func run() error {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	var req plugin.Request
	err = proto.Unmarshal(b, &req)
	if err != nil {
		return err
	}

	var response proto.Message

	switch req := req.Req.(type) {
	case *plugin.Request_VersionRequest:
		response = &plugin.VersionResponse{Version: "1.0.0"}

	case *plugin.Request_SyntaxRequest:
		var buf bytes.Buffer
		err := quick.Highlight(&buf, req.SyntaxRequest.Code, req.SyntaxRequest.Language, "html", "monokai")
		if err != nil {
			return err
		}
		response = &plugin.SyntaxResponse{Output: buf.String()}
	}

	b, err = proto.Marshal(response)
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(b)
	return err
}

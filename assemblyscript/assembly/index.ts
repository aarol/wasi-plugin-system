// The entry file of your WebAssembly module.
import { decodeRequest } from '../wasi_plugin/Request'
import { VersionResponse, encodeVersionResponse } from '../wasi_plugin/VersionResponse'

const input = new ArrayBuffer(2048)
process.stdin.read(input)

const req = decodeRequest(Uint8Array.wrap(input))

if (req.versionRequest) {
  const res = encodeVersionResponse(new VersionResponse("1.0.0"))
  process.stdout.write(String.UTF8.decode(res.buffer))
}
// Code generated by protoc-gen-as. DO NOT EDIT.
// Versions:
//   protoc-gen-as v1.3.0
//   protoc        v5.26.0--rc1

import { Writer, Reader, Protobuf } from "as-proto/assembly";

export class VersionResponse {
  static encode(message: VersionResponse, writer: Writer): void {
    writer.uint32(10);
    writer.string(message.version);
  }

  static decode(reader: Reader, length: i32): VersionResponse {
    const end: usize = length < 0 ? reader.end : reader.ptr + length;
    const message = new VersionResponse();

    while (reader.ptr < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.version = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  }

  version: string;

  constructor(version: string = "") {
    this.version = version;
  }
}

export function encodeVersionResponse(message: VersionResponse): Uint8Array {
  return Protobuf.encode(message, VersionResponse.encode);
}

export function decodeVersionResponse(buffer: Uint8Array): VersionResponse {
  return Protobuf.decode<VersionResponse>(buffer, VersionResponse.decode);
}

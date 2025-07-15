# goc

goc is the abbreviation for go component, which is a lightweight framework.

## features

* database: gorm
* redis: go-redis
* log: zap
* config: yaml, viper
* auth: jwt
* http: gin
* rpc: grpc, grpc-gateway

## Getting Started

### Prerequisites

### wire

<https://github.com/google/wire>

```bash
go install github.com/google/wire/cmd/wire@latest
```

### stringer

<https://pkg.go.dev/golang.org/x/tools/cmd/stringer>

```bash
go install golang.org/x/tools/cmd/stringer@latest
```

### buf

<https://buf.build/docs/installation/>

```bash
go install github.com/bufbuild/buf/cmd/buf@latest
```

## Docs

Protocol Buffers Documentation: <https://protobuf.dev/>

gRPC-Gateway: <https://grpc-ecosystem.github.io/grpc-gateway/>

Transcoding HTTP/JSON to gRPC: <https://cloud.google.com/endpoints/docs/grpc/transcoding>

Buf: <https://buf.build/docs/introduction>

protovalidate: <https://github.com/bufbuild/protovalidate>

## Style guide

API design guide: <https://cloud.google.com/apis/design>

Programming Guides: <https://protobuf.dev/programming-guides/>

Style guide: <https://buf.build/docs/best-practices/style-guide>

Uber Protobuf Style Guide V2: <https://github.com/uber/prototool/blob/dev/style/README.md>

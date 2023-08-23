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

## doc

https://grpc-ecosystem.github.io/grpc-gateway/

https://cloud.google.com/endpoints/docs/grpc/transcoding?hl=zh-cn

https://github.com/bufbuild/protovalidate


## use protoc generate code

### protoc

<https://github.com/protocolbuffers/protobuf>

```bash
cd third_party
git clone https://github.com/googleapis/googleapis.git --depth 1
git clone https://github.com/bufbuild/protovalidate.git --depth 1
```

### plugins

grpc-gateway: <https://grpc-ecosystem.github.io/grpc-gateway/docs/tutorials/introduction>

```bash
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/envoyproxy/protoc-gen-validate@latest
```

This will place four binaries in your $GOBIN;

```bash
protoc -I ./component/auth/api/v1 \
		-I ./third_party/googleapis \
		-I ./third_party/protovalidate/proto/protovalidate \
		--go_out ./component/auth/api/v1 --go_opt paths=source_relative \
		--go-grpc_out ./component/auth/api/v1 --go-grpc_opt paths=source_relative \
		--grpc-gateway_out ./component/auth/api/v1 --grpc-gateway_opt paths=source_relative \
		--validate_out=lang=go,paths=source_relative:./component/auth/api/v1 \
		--openapiv2_out ./component/auth/api/v1 --openapiv2_opt logtostderr=true --openapiv2_opt use_go_templates=true \
		./component/auth/api/v1/auth.proto
```

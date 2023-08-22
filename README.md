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


## grpc

```bash
cd third_party
git clone https://github.com/googleapis/googleapis.git --depth 1
```

```bash
	protoc -I ./service/auth/api/v1 \
		-I ./third_party/googleapis \
		--go_out ./service/auth/api/v1 --go_opt paths=source_relative \
		--go-grpc_out ./service/auth/api/v1 --go-grpc_opt paths=source_relative \
		--grpc-gateway_out ./service/auth/api/v1 --grpc-gateway_opt paths=source_relative \
		./service/auth/api/v1/auth.proto
	protoc -I ./service/file/api/v1 \
		-I ./third_party/googleapis \
		--go_out ./service/file/api/v1 --go_opt paths=source_relative \
		--go-grpc_out ./service/file/api/v1 --go-grpc_opt paths=source_relative \
		--grpc-gateway_out ./service/file/api/v1 --grpc-gateway_opt paths=source_relative \
		./service/file/api/v1/file.proto
```

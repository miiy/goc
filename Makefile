TARGET=server

.PHONY: proto
proto:
	@protoc -I ./component/auth/api/v1 \
		-I ./third_party/googleapis \
		--go_out ./component/auth/api/v1 --go_opt paths=source_relative \
		--go-grpc_out ./component/auth/api/v1 --go-grpc_opt paths=source_relative \
		--grpc-gateway_out ./component/auth/api/v1 --grpc-gateway_opt paths=source_relative \
		./component/auth/api/v1/auth.proto
	@protoc -I ./component/file/api/v1 \
		-I ./third_party/googleapis \
		--go_out ./component/file/api/v1 --go_opt paths=source_relative \
		--go-grpc_out ./component/file/api/v1 --go-grpc_opt paths=source_relative \
		--grpc-gateway_out ./component/file/api/v1 --grpc-gateway_opt paths=source_relative \
		./component/file/api/v1/file.proto

	@protoc -I ./examples/apiserver/api/echo/v1 \
		-I ./third_party/googleapis \
		--go_out ./examples/apiserver/api/echo/v1 --go_opt paths=source_relative \
		--go-grpc_out ./examples/apiserver/api/echo/v1 --go-grpc_opt paths=source_relative \
		--grpc-gateway_out ./examples/apiserver/api/echo/v1 --grpc-gateway_opt paths=source_relative \
		./examples/apiserver/api/echo/v1/echo.proto


.PHONY: help
help:
	@echo "make proto: proto file"

.PHONY: tools
tools:
	# go mod tidy
	go get \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
		github.com/golang/protobuf/protoc-gen-go \
		github.com/rakyll/statik
	# git clone https://github.com/grpc-ecosystem/grpc-gateway

build:
	@echo Compiling proto files
	@protoc \
		-I./example \
		-I./third_party/grpc-gateway/ \
		-I./third_party/googleapis \
		--go_out=plugins=grpc,paths=source_relative:./example \
		--grpc-gateway_out=./example \
		--swagger_out=./third_party/OpenAPI/ \
		./example/your_service.proto

	@statik -m -f -src third_party/OpenAPI/

	@echo Building example-grpc-gateway
	@go build -v -o example-grpc-gateway github.com/pieterclaerhout/example-grpc-gateway

run-server: build
	./example-grpc-gateway --what=server
	
run-grpc-client: build
	./example-grpc-gateway --what=grpc-client
	
run-rest-client: build
	./example-grpc-gateway --what=rest-client
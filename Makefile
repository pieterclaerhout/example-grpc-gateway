.PHONY: tools
tools:
	go mod tidy
	go install \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
		github.com/golang/protobuf/protoc-gen-go
	git clone https://github.com/grpc-ecosystem/grpc-gateway

build:
	@echo Compiling proto files
	@protoc -I. \
		-I./grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:./example \
		--grpc-gateway_out=logtostderr=true:. \
		--swagger_out=logtostderr=true:. \
		./example/your_service.proto
	@#protoc -I/usr/local/include -I. \
	@# 	-I$GOPATH/src \
	@# 	-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	@# 	--grpc-gateway_out=logtostderr=true:. \
	@# 	path/to/your_service.proto
	@echo Building example-grpc-gateway
	@go build -v -o example-grpc-gateway github.com/pieterclaerhout/example-grpc-gateway

run-server: build
	./example-grpc-gateway --what=server
	
run-client: build
	./example-grpc-gateway --what=client
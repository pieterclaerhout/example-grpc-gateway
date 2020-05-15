PROJECT_NAME := example-grpc-gateway
REVISION := 1.0
DOCKER_ACCOUNT := pieterclaerhout

install-tools:
	@go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	@go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	@go get github.com/golang/protobuf/protoc-gen-go
	@go get github.com/rakyll/statik

generate:
	@echo Compiling proto files
	@protoc \
		-I./example \
		-I./third_party/grpc-gateway/ \
		-I./third_party/googleapis \
		--go_out=plugins=grpc,paths=source_relative:./example \
		--grpc-gateway_out=./example \
		--swagger_out=./third_party/OpenAPI/ \
		./example/your_service.proto
	@echo Packaging resources
	@statik -m -f -src third_party/OpenAPI/

build: generate
	@echo Building $(PROJECT_NAME)
	@go build -v -o $(PROJECT_NAME) github.com/pieterclaerhout/$(PROJECT_NAME)

publish: build
	@docker build -t $(PROJECT_NAME) .
	@docker tag $(PROJECT_NAME) $(DOCKER_ACCOUNT)/$(PROJECT_NAME):$(REVISION)
	@docker push $(DOCKER_ACCOUNT)/$(PROJECT_NAME):$(REVISION)

run-server: build
	@./$(PROJECT_NAME) --what=server
	
run-grpc-client: build
	@./$(PROJECT_NAME) --what=grpc-client
	
run-rest-client: build
	@./$(PROJECT_NAME) --what=rest-client

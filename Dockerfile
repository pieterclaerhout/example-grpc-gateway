# Step 1 - Create layer with the downloaded modules
FROM golang:1.14.2-alpine AS mod-download

ENV PROTOC_VERSION 3.11.4

RUN apk --no-cache add git make unzip protobuf protoc

RUN mkdir -p /app

ADD Makefile /app
ADD go.mod /app
ADD go.sum /app

WORKDIR /app

RUN wget https://github.com/google/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip \
    && unzip protoc-${PROTOC_VERSION}-linux-x86_64.zip -x readme.txt -d /usr/local \
    && rm /usr/local/bin/protoc

RUN make install-tools
RUN go mod download


# Step 2 - Build
FROM mod-download AS builder

ADD . /app
WORKDIR /app

RUN PATH="/usr/local/bin:${PATH}" make build
RUN chmod a+x /app/example-grpc-gateway


# Step 3 - Final
FROM alpine:3.11 

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/example-grpc-gateway /

ENTRYPOINT ["/example-grpc-gateway"]
EXPOSE 8080

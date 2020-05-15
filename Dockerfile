# Step 1 - Create layer with the downloaded modules
FROM golang:1.14.2-alpine AS mod-download

RUN mkdir -p /app

ADD go.mod /app
ADD go.sum /app

WORKDIR /app

RUN go mod download


# Step 2 - Build
FROM mod-download AS builder

ADD . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath --ldflags '-extldflags -static' -o example-environ-server github.com/twixlmedia/example-environ-server


# Step 3 - Final
FROM alpine:3.11 

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/example-environ-server /
RUN chmod a+x /example-environ-server

ENV TWX_ENVIRONMENT=production

ENTRYPOINT ["/example-environ-server"]
EXPOSE 8080

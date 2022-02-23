all: test lint vet build

build: client server csrparse

client:
	@cd cmd/$@ && go build -o ../../bin/$@

server:
	@cd cmd/$@ && go build -o ../../bin/$@

csrparse:
	@cd cmd/$@ && go build -o ../../bin/$@

test:
	@go test ./...

lint:
	@revive ./...

vet:
	@go vet ./...

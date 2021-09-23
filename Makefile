
.PHONY:build
build: gen-examples
	go build -o ./bin/goqu-crud-gen ./cmd/goqu-crud-gen


.PHONY:install
install:
	go install ./cmd/...

.PNONY:gen-examples
gen-examples:
	cd ./examples/example1/ && go generate ./...
	cd ./examples/example2/ && go generate ./...
	cd ./examples/example3/ && go generate ./...


.PHONY:build
build:
	go build -o ./bin/goqu-crud-gen ./cmd/goqu-crud-gen


.PHONY:install
install:
	go install ./cmd/...

.PNONY:gen-examples
gen-examples:
	cd ./examples/example1/ && go generate ./... && go vet ./...
	cd ./examples/example2/ && go generate ./... && go vet ./...
	cd ./examples/example3/ && go generate ./... && go vet ./...

.PHONY:before-commit
before-commit: build install gen-examples
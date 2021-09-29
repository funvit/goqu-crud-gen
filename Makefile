
.PHONY:build
build:
	go build -o ./bin/goqu-crud-gen ./cmd/goqu-crud-gen


.PHONY:install
install:
	go install ./cmd/...


.PHONY:generate-examples
generate-examples:
	go generate ./examples/example1/model/user.go
	go generate ./examples/example2/user.go
	go generate ./examples/example3/adapters/mysql/account.go
	go generate ./examples/example3/adapters/mysql/user_public_fields.go
	go vet ./examples/...
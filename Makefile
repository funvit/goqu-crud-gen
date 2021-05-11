
.PHONY:build
build:
	go build -o ./bin/goqu-crud-gen ./cmd/goqu-crud-gen


.PHONY:install
install:
	go install ./cmd/...
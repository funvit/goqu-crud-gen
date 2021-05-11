
.PHONY:run
run:
	docker run  --name example3 \
		-e MYSQL_ROOT_PASSWORD=root_pass \
		-e MYSQL_USER=test_login \
		-e MYSQL_PASSWORD=test_pass \
		-e MYSQL_DATABASE=example3 \
		-p 127.0.0.1:53306:3306 \
		-d \
		mysql:5.7

	@echo Mysql needs some time to boot...
	sleep 10

	@echo ""
	go run ./cmd/...

	@echo ""
	@echo Cleanup...
	docker stop example3
	docker rm example3

.PHONY:generate
generate:
	go generate ./...



run:
	go run cmd/analyzer/main.go -file testdata/logs/nginx.log -config config/config.yaml -type nginx

test:
	go test -v -coverprofile coverage.out ./...

test-cover: test
	go tool cover -html coverage.out -o coverage.html

helper:
	go run helper/helper.go

.PHONY: run, test, test-cover helper

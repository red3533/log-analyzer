run:
	go run cmd/analyzer/main.go -file testdata/logs/nginx.log -config config/config.yaml

helper:
	go run helper/helper.go

.PHONY: run, helper

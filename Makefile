.PHONY: mvx

build:
	@go build -o bin/api
run:
	go run ./main.go serve
mvx:
	@go run test_mvx/test.go
fixtures:
	@go run scripts/fixtures.go serve
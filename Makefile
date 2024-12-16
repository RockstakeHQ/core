build:
	@go build -o bin/api
run:
	@./bin/api
mvx:
	@go run ./main.go
fixtures:
	@go run scripts/fixtures.go serve
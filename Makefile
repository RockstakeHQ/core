build:
	@go build -o bin/api
run:
	go run ./main.go serve
mvx:
	@go run ./main.go 
fixtures:
	@go run scripts/fixtures.go serve
build:
	@go build -o bin/api
run:
	@./bin/api
seed:
	@go run scripts/seed.go
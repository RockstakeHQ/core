build:
	@go build -o bin/api
	
run:
	@./bin/api

get_sportsbook:
	@go run services/data_processing_services/sportsbook_aggregation_service/main/main.go

run_sportsbook:
	@go build -o bin/sportsbook_service /Users/andrewkhirita/Desktop/betcube_engine/services/sportsbook_service/main.go
	@./bin/sportsbook_service

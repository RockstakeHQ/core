package main

import (
	"betcube-engine/services/data_processing_services/types_api"
	"betcube-engine/types"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getJSON(url string, params map[string]string, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	oddsEndpoint := "https://v3.football.api-sports.io/odds"
	fixturesEndpoint := "https://v3.football.api-sports.io/fixtures"
	date := "2024-09-13"
	bookmaker := 8

	headers := map[string]string{
		"x-rapidapi-host": "v3.football.api-sports.io",
		"x-rapidapi-key":  "33db0d9d9531386613988c43458700a6",
	}

	fixturesParams := map[string]string{
		"date": date,
	}
	fixturesData, err := getJSON(fixturesEndpoint, fixturesParams, headers)
	if err != nil {
		log.Fatalf("Error fetching fixtures: %v", err)
	}

	var fixtures types_api.FixtureResponse
	if err := json.Unmarshal(fixturesData, &fixtures); err != nil {
		log.Fatalf("Error parsing fixtures JSON: %v", err)
	}

	var allOddsData []types_api.Odds
	oddsParams := map[string]string{
		"date":      date,
		"bookmaker": strconv.Itoa(bookmaker),
	}

	for {
		oddsData, err := getJSON(oddsEndpoint, oddsParams, headers)
		if err != nil {
			log.Fatalf("Error fetching odds: %v", err)
		}

		var oddsPage types_api.OddsResponse
		if err := json.Unmarshal(oddsData, &oddsPage); err != nil {
			log.Fatalf("Error parsing odds JSON: %v", err)
		}

		allOddsData = append(allOddsData, oddsPage.Response...)
		if oddsPage.Paging.Current == oddsPage.Paging.Total {
			break
		}
		oddsParams["page"] = strconv.Itoa(oddsPage.Paging.Current + 1)
	}

	var matchInfo []types_api.MatchData
	for _, fixture := range fixtures.Response {
		fixtureID := fixture.Fixture.ID
		var matchOdds []types_api.Odds
		for _, odd := range allOddsData {
			if odd.Fixture.ID == fixtureID {
				matchOdds = append(matchOdds, odd)
			}
		}

		updateValue := ""
		if len(matchOdds) > 0 {
			updateValue = matchOdds[0].Update
		}

		matchData := types_api.MatchData{
			Fixture: fixture.Fixture,
			League:  fixture.League,
			Teams:   fixture.Teams,
			Goals:   fixture.Goals,
			Score:   fixture.Score,
			Update:  updateValue,
		}

		if len(matchOdds) > 0 {
			for _, bookmaker := range matchOdds[0].Bookmakers {
				for _, bet := range bookmaker.Bets {
					betData := types_api.BetData{
						ID:   bet.ID,
						Name: bet.Name,
					}
					for _, value := range bet.Values {
						valueData := types_api.ValueData{
							Value: value.Value,
							Odd:   value.Odd,
						}
						betData.Values = append(betData.Values, valueData)
					}
					matchData.Bets = append(matchData.Bets, betData)
				}
			}
		}

		matchInfo = append(matchInfo, matchData)
	}

	file, err := os.Create("match_info.json")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(matchInfo); err != nil {
		log.Fatalf("Error encoding JSON to file: %v", err)
	}

	if len(matchInfo) > 0 {
		log.Printf("Succeses!")
	}

	//part 2

	mongoEndpoint := os.Getenv("MONGO_DB_URL")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	betTypesCollection := client.Database(dbName).Collection("bet_types")
	betValuesCollection := client.Database(dbName).Collection("bet_values")
	fixturesCollection := client.Database(dbName).Collection("fixtures")
	leaguesCollection := client.Database(dbName).Collection("leagues")
	teamsCollection := client.Database(dbName).Collection("teams")
	sportsbooksCollection := client.Database(dbName).Collection("sportsbooks")

	fileContent, err := os.ReadFile("match_info.json")
	if err != nil {
		log.Fatal("Error read JSON file:", err)
	}

	var matches []types.Match
	err = json.Unmarshal(fileContent, &matches)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
	}

	var updatedMatches []types.Match
	for _, match := range matches {
		var newBets []types.Bet
		for _, bet := range match.Bets {
			var betType types.BetType
			err := betTypesCollection.FindOne(context.Background(), bson.M{"_id": bet.ID}).Decode(&betType)
			if err != nil {
				log.Printf("Nu s-a găsit tipul de pariu cu ID %d: %v", bet.ID, err)
				continue
			}

			bet.Name = betType.Name

			var newValues []types.Value
			for _, betValue := range bet.Values {
				// Căutăm valoarea în baza de date
				var dbBetValue types.BetValue
				err := betValuesCollection.FindOne(context.Background(), bson.M{"value": betValue.Value}).Decode(&dbBetValue)
				if err != nil {
					log.Printf("Nu s-a găsit valoarea pariului: %v", err)
					continue
				}

				// Actualizarea valorii pariului
				betValue.ID = dbBetValue.ID
				newValues = append(newValues, betValue)
			}

			if len(newValues) > 0 {
				bet.Values = newValues
				newBets = append(newBets, bet)
			}
		}
		if len(newBets) > 0 {
			match.Bets = newBets
			updatedMatches = append(updatedMatches, match)
		}
	}

	outputFileContent, err := json.MarshalIndent(updatedMatches, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("updated_match_info.json", outputFileContent, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//part 3

	jsonData, err := os.ReadFile("updated_match_info.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	var matchesDb []types_api.MatchDb
	err = json.Unmarshal(jsonData, &matchesDb)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	for _, match := range matchesDb {
		flagURL := match.League.Flag
		if match.League.Country == "World" {
			flagURL = "https://upload.wikimedia.org/wikipedia/en/6/6b/Terrestrial_globe.svg"
		}

		league := types.League{
			ID:      match.League.ID,
			Name:    match.League.Name,
			Country: match.League.Country,
			Logo:    match.League.Logo,
			Flag:    flagURL,
		}

		_, err := leaguesCollection.UpdateOne(
			context.TODO(),
			map[string]interface{}{"_id": league.ID},
			map[string]interface{}{"$set": league},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			fmt.Println("Error inserting/updating league:", err)
		} else {
			fmt.Println("Upserted league with ID:", league.ID)
		}

		// Inserează echipele Home și Away
		homeTeam := types.Team{
			ID:     match.Teams.Home.ID,
			Name:   match.Teams.Home.Name,
			Logo:   match.Teams.Home.Logo,
			Winner: match.Teams.Home.Winner,
		}

		_, err = teamsCollection.UpdateOne(
			context.TODO(),
			map[string]interface{}{"_id": homeTeam.ID},
			map[string]interface{}{"$set": homeTeam},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			fmt.Println("Error inserting/updating home team:", err)
		} else {
			fmt.Println("Upserted home team with ID:", homeTeam.ID)
		}

		awayTeam := types.Team{
			ID:     match.Teams.Away.ID,
			Name:   match.Teams.Away.Name,
			Logo:   match.Teams.Away.Logo,
			Winner: match.Teams.Away.Winner,
		}

		_, err = teamsCollection.UpdateOne(
			context.TODO(),
			map[string]interface{}{"_id": awayTeam.ID},
			map[string]interface{}{"$set": awayTeam},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			fmt.Println("Error inserting/updating away team:", err)
		} else {
			fmt.Println("Upserted away team with ID:", awayTeam.ID)
		}

		if match.Update != "" && len(match.Bets) > 0 {
			var updateTime time.Time
			if match.Update != "" {
				updateTime, err = time.Parse(time.RFC3339, match.Update)
				if err != nil {
					fmt.Println("Error parsing update time:", err)
					continue
				}
			} else {
				updateTime = time.Time{}
			}

			sportsbook := types.Sportsbook{
				ID:     match.Fixture.ID,
				Update: updateTime,
				Bets:   match.Bets,
			}

			// Verifică dacă sportsbook există
			var existingSportsbook types.Sportsbook
			err := sportsbooksCollection.FindOne(context.TODO(), map[string]interface{}{"_id": sportsbook.ID}).Decode(&existingSportsbook)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					// Nu există, deci inserăm
					_, err = sportsbooksCollection.InsertOne(context.TODO(), sportsbook)
					if err != nil {
						fmt.Println("Error inserting sportsbook:", err)
					} else {
						fmt.Println("Inserted sportsbook with ID:", sportsbook.ID)
					}
				} else {
					// Eroare la găsirea documentului
					fmt.Println("Error finding sportsbook:", err)
				}
			} else {
				// Există, deci actualizăm
				_, err = sportsbooksCollection.UpdateOne(
					context.TODO(),
					map[string]interface{}{"_id": sportsbook.ID},
					map[string]interface{}{"$set": sportsbook},
				)
				if err != nil {
					fmt.Println("Error updating sportsbook:", err)
				} else {
					fmt.Println("Updated sportsbook with ID:", sportsbook.ID)
				}
			}

			fixture := types.Fixture{
				ID:           match.Fixture.ID,
				Date:         match.Fixture.Date,
				League:       match.League.ID,
				Teams:        []int{match.Teams.Home.ID, match.Teams.Away.ID},
				Status:       match.Fixture.Status,
				Score:        match.Score,
				Goals:        match.Goals,
				SportsbookID: sportsbook.ID,
			}

			// Verifică dacă fixture există
			var existingFixture types.Fixture
			err = fixturesCollection.FindOne(context.TODO(), map[string]interface{}{"_id": fixture.ID}).Decode(&existingFixture)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					// Nu există, deci inserăm
					_, err = fixturesCollection.InsertOne(context.TODO(), fixture)
					if err != nil {
						fmt.Println("Error inserting fixture:", err)
					} else {
						fmt.Println("Inserted fixture with ID:", fixture.ID)
					}
				} else {
					// Eroare la găsirea documentului
					fmt.Println("Error finding fixture:", err)
				}
			} else {
				// Există, deci actualizăm
				_, err = fixturesCollection.UpdateOne(
					context.TODO(),
					map[string]interface{}{"_id": fixture.ID},
					map[string]interface{}{"$set": fixture},
				)
				if err != nil {
					fmt.Println("Error updating fixture:", err)
				} else {
					fmt.Println("Updated fixture with ID:", fixture.ID)
				}
			}
		}
	}
}

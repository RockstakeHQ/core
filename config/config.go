package config

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	MongoClient *mongo.Client
	PostgresDB  *bun.DB
	StripeKey   string
}

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
}

func InitMongoClient() *mongo.Client {
	mongoEndpoint := os.Getenv("MONGO_DB_URL")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	return client
}

func InitPostgresDB() *bun.DB {
	supabaseEndpoint := os.Getenv("SUPABASE_URL")
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(supabaseEndpoint)))
	return bun.NewDB(sqldb, pgdialect.New())
}

func InitStripeKey() string {
	stripeKey := os.Getenv("STRIPE_API_SECRET")
	if stripeKey == "" {
		log.Fatal("Stripe API Secret not found in environment")
	}
	return stripeKey
}

func NewConfig() *Config {
	LoadEnv()

	return &Config{
		MongoClient: InitMongoClient(),
		PostgresDB:  InitPostgresDB(),
		StripeKey:   InitStripeKey(),
	}
}

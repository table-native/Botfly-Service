package db

import (
	"context"
	"os"
	"time"

	"github.com/table-native/Botfly-Service/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type BotflyDb struct {
	Db *mongo.Database
}

func NewBotflyDb() *BotflyDb {
	mongo_uri := os.Getenv("MONGO_CONN")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongo_uri))
	if err != nil {
		logger.Fatal("Failed to connect to mongo", zap.Error(err))
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		logger.Fatal("Failed to connect to mongo", zap.Error(err))
	}

	db := client.Database("BotFly")

	return &BotflyDb{
		Db: db,
	}
}

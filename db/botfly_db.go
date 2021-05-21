package db

import (
	"context"
	"crypto/tls"
	"os"
	"strings"
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
	mongo_uri := strings.TrimSuffix(os.Getenv("MONGO_CONN"), "\n")

	mongoOpts := options.Client().ApplyURI(mongo_uri)
	mongoOpts.TLSConfig.MinVersion = tls.VersionTLS12
	mongoOpts.TLSConfig.InsecureSkipVerify = true

	client, err := mongo.NewClient(mongoOpts)
	if err != nil {
		logger.Fatal("Failed to connect to mongo", zap.Error(err))
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		logger.Fatal("Failed to connect to mongo", zap.Error(err))
	}

	err = client.Ping(nil, nil)
	if err != nil {
		logger.Fatal("Failed to connect to mongo", zap.Error(err))
	}

	db := client.Database("BotFly")

	return &BotflyDb{
		Db: db,
	}
}

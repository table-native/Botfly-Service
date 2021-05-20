package db

import (
	"context"
	"encoding/base64"

	"github.com/table-native/Botfly-Service/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type ScriptsDto struct {
	collection *mongo.Collection
}

type ScriptsModel struct {
	Id           string `bson:"_id"`
	UserId       string `bson:"userId"`
	UserCode     string `bson:"userCode"`
	PlatformCode string `bson:"platformCode"`
	Game         string `bson:"game"`
}

func (s ScriptsModel) GetId() string {
	base64Id := base64.StdEncoding.EncodeToString([]byte(s.UserId + "/" + s.Game))
	return base64Id
}

func NewScriptsDto(db *BotflyDb) *ScriptsDto {
	return &ScriptsDto{
		collection: db.Db.Collection("Scripts"),
	}
}

func (s ScriptsDto) Create(scriptsModel ScriptsModel) chan ScriptsModel {
	ch := make(chan ScriptsModel)

	go func() {
		scriptsModel.Id = scriptsModel.GetId()

		document, _ := bson.Marshal(scriptsModel)
		_, err := s.collection.InsertOne(context.Background(), document)

		if err != nil {
			logger.Error("Failed inserting into DB", zap.Error(err))
		}

		ch <- scriptsModel
	}()

	return ch
}

func (s ScriptsDto) FindByUserAndGame(userId string, game string) chan ScriptsModel {
	ch := make(chan ScriptsModel)

	go func() {
		scriptsModel := ScriptsModel{UserId: userId, Game: game}
		id := scriptsModel.GetId()

		logger.Info("Query id is", zap.Any("id", id))
		scriptsBson := s.collection.FindOne(context.Background(), bson.M{"_id": id})

		if scriptsBson.Err() != nil {
			logger.Error("Got error fetching user script", zap.Error(scriptsBson.Err()))
		}

		scripts := &ScriptsModel{}
		scriptsBson.Decode(scripts)

		ch <- *scripts
	}()
	return ch
}

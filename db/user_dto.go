package db

import (
	"context"
	"encoding/base64"

	"github.com/table-native/Botfly-Service/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type UserDto struct {
	collection *mongo.Collection
}

type UserModel struct {
	Id    string `bson:"_id"`
	Email string `bson:"email"`
	Score int    `bson:"score"`
}

func (u UserModel) GetId() string {
	uid := base64.StdEncoding.EncodeToString([]byte(u.Email))
	return uid
}

func NewUserDto(db *BotflyDb) *UserDto {
	return &UserDto{
		collection: db.Db.Collection("Users"),
	}
}

func (u UserDto) Create(user UserModel) chan UserModel {
	ch := make(chan UserModel)

	go func() {
		user.Id = user.GetId()

		document, _ := bson.Marshal(user)
		_, err := u.collection.InsertOne(context.Background(), document)

		if err != nil {
			logger.Error("Failed inserting into DB", zap.Error(err))
		}

		ch <- user
	}()

	return ch
}

func (u UserDto) FindUserById(id string) chan UserModel {
	ch := make(chan UserModel)

	go func() {
		userBson := u.collection.FindOne(context.Background(), bson.M{"_id": id})

		user := UserModel{}
		userBson.Decode(user)
		ch <- user
	}()
	return ch
}

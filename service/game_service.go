package service

import (
	"context"
	"os"
	"path"
	"text/template"

	"github.com/table-native/Botfly-Service/auth"
	"github.com/table-native/Botfly-Service/games"
	pb "github.com/table-native/Botfly-Service/generated"
	"github.com/table-native/Botfly-Service/logger"
	"go.uber.org/zap"
)

type GameService struct {
	pb.UnimplementedGameServiceServer
}

func (g GameService) GetBotTemplate(ctx context.Context, gameDetails *pb.GameDetails) (*pb.BotTemplate, error) {
	return &pb.BotTemplate{
		Template: games.Tic_Tac_Toe_Template,
	}, nil
}

func saveUserCode(userId string, code string) {
	userCodeTemplate := template.Must(
		template.New("userCode").Parse(games.Tic_Tac_Toe_Driver_Code),
	)

	os.MkdirAll(path.Join("data", userId), os.ModePerm)
	codeFile, err := os.Create(path.Join("data", userId, "ticTacToe.js"))
	if err != nil {
		logger.Fatal("Failed writing usercode file", zap.Error(err))
		return
	}
	defer codeFile.Close()

	//delete current file
	codeFile.Truncate(0)
	err = userCodeTemplate.Execute(codeFile, code)
}

func (g GameService) SaveMyBot(ctx context.Context, code *pb.BotTemplate) (*pb.SaveStatus, error) {
	userId := ctx.Value(auth.Claims("userId")).(string)
	logger.Info("Saving bot for userId", zap.String("userId", userId))

	saveUserCode(userId, code.Template)
	return &pb.SaveStatus{}, nil
}

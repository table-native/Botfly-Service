package service

import (
	"bytes"
	"context"
	"text/template"

	"github.com/table-native/Botfly-Service/auth"
	"github.com/table-native/Botfly-Service/db"
	"github.com/table-native/Botfly-Service/games"
	pb "github.com/table-native/Botfly-Service/generated"
	"github.com/table-native/Botfly-Service/logger"
	"go.uber.org/zap"
)

type GameService struct {
	pb.UnimplementedGameServiceServer
	scriptsDto *db.ScriptsDto
}

func NewGameService(scriptsDto *db.ScriptsDto) *GameService {
	return &GameService{
		scriptsDto: scriptsDto,
	}
}

func (g GameService) GetBotTemplate(ctx context.Context, gameDetails *pb.GameDetails) (*pb.BotTemplate, error) {
	return &pb.BotTemplate{
		Template: games.Tic_Tac_Toe_Template,
	}, nil
}

func getUserCodeSubstitutingMain(code string) string {
	userCodeTemplate := template.Must(
		template.New("userCode").Parse(games.Tic_Tac_Toe_Driver_Code),
	)

	var outputCode bytes.Buffer
	userCodeTemplate.Execute(&outputCode, code)
	return outputCode.String()
}

func (g GameService) SaveMyBot(ctx context.Context, code *pb.BotTemplate) (*pb.SaveStatus, error) {
	userId := ctx.Value(auth.Claims("userId")).(string)
	logger.Info("Saving bot for userId", zap.String("userId", userId))

	srcCode := getUserCodeSubstitutingMain(code.Template)
	scriptsModel := db.ScriptsModel{
		UserId:       userId,
		PlatformCode: srcCode,
		UserCode:     code.Template,
		Game:         "tic-tac-toe",
	}

	<-g.scriptsDto.Create(scriptsModel)
	return &pb.SaveStatus{}, nil
}

func (g GameService) GetMyBot(ctx context.Context, _ *pb.EmptyRequest) (*pb.BotTemplate, error) {
	userId := ctx.Value(auth.Claims("userId")).(string)
	scritpsModel := <-g.scriptsDto.FindByUserAndGame(userId, "tic-tac-toe")
	return &pb.BotTemplate{Template: scritpsModel.UserCode}, nil
}

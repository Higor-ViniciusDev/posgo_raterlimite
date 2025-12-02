package main

import (
	"fmt"
	"os"

	"github.com/Higor-ViniciusDev/posgo_raterlimite/configuration/database"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/configuration/logger"
	db_infra "github.com/Higor-ViniciusDev/posgo_raterlimite/internal/infra/database"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/infra/web/controller"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/infra/web/middleware"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/infra/web/server"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/usecase/expire_usecase"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/usecase/policy_usecase"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/usecase/tolken_usecase"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	defer logger.GetLogger().Sync()
	// ctx := context.Background()

	if err := godotenv.Load("cmd/raterlimite/.env"); err != nil {
		logger.Error("Erro ao carregar variaveis de ambiente", err)
		return
	}
	redis := database.NewConnectionRedis()
	tolkenController, policyUsecase := initDependeces(redis)

	webServerPort := os.Getenv("WEB_SERVER_POR")
	webServer := server.NovoWebServer(fmt.Sprintf(":%v", webServerPort))
	webServer.RegistrarRota("/tolken", tolkenController.CreateTolken, "POST")
	webServer.RegistrarRota("/", nil, "GET", middleware.RateLimiterMiddleware(&policyUsecase))

	webServer.IniciarWebServer()
}

func initDependeces(redisCli *redis.Client) (controller.TolkenController, policy_usecase.PolicyUsecase) {
	var tolkenController controller.TolkenController

	expirerUsecase := expire_usecase.NewDefaultExpirer()

	//Tolken dependeces
	tolkeRepository := db_infra.NewTolkenDB(redisCli)
	tolkenUsecase := tolken_usecase.NewTolkenUsecase(tolkeRepository, expirerUsecase)
	tolkenController = *controller.NewTolkenController(tolkenUsecase)

	requestRespository := db_infra.NewRequestInfoRepository(redisCli)

	policyUsecase := *policy_usecase.NewPolicyUsecase(expirerUsecase, tolkeRepository, requestRespository)

	return tolkenController, policyUsecase
}

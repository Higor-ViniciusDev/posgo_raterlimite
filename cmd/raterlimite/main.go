package main

import (
	"github.com/Higor-ViniciusDev/posgo_raterlimite/configuration/logger"
	"github.com/joho/godotenv"
)

func main() {
	defer logger.GetLogger().Sync()
	// ctx := context.Background()

	if err := godotenv.Load("cmd/raterlimite/.env"); err != nil {
		logger.Error("Erro ao carregar variaveis de ambiente", err)
		return
	}

	// Output: allowed 1 remaining 9
}

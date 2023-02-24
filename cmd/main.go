package main

import (
	"fmt"
	"os"

	"go.uber.org/zap" // logger

	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/delivery"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/repository"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/usecase"
	"github.com/go-park-mail-ru/2023_1_Technokaif/server"
)

// Here for a while
const PORT = "4443"

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Logger can not be defined: "+err.Error())
	}

	repository := repository.NewRepository()
	services := usecase.NewUsecase(repository)
	handler := delivery.NewHandler(services)

	server := new(server.Server)
	if err := server.Run(PORT, handler.InitRouter()); err != nil {
		logger.Error("Error while launching server: " + err.Error())
	}
}

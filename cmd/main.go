package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2023_1_Technokaif/fluire"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/delivery"
)

// Here for a while
const PORT = "8080"

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Logger can not be defined: "+err.Error())
	}

	router := delivery.InitRouter()

	server := new(fluire.Server)
	if err := server.Run(PORT, router); err != nil {
		logger.Error("Error while launching server: " + err.Error())
	}
}

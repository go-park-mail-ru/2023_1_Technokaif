package main

import (
	"fmt"
	"os"

	"go.uber.org/zap" // logger

	"github.com/go-park-mail-ru/2023_1_Technokaif/server"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/delivery"
)

// Here for a while
const PORT = "4443"

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Logger can not be defined: " + err.Error())
	}

	handler := new(delivery.Handler)

	server := new(server.Server)
	if err := server.Run(PORT, handler.InitRouter()); err != nil {
		logger.Error("Error while launching server: " + err.Error())
	}
}

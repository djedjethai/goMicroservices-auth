package main

import (
	"github.com/djedjethai/bankingAuth/app"
	"github.com/djedjethai/bankingAuth/logger"
)

func main() {
	logger.Info("start the application")
	app.Start()

}

package main

import (
	"github.com/djedjethai/bankingAuth/app"
	"github.com/djedjethai/bankingLib/logger"
)

func main() {
	logger.Info("start the application")
	app.Start()

}

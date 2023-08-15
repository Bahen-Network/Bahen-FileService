package main

import (
	"file-service/config"
	"file-service/objectstorage"
	"file-service/router"
	common "file-service/util"
)

// main.go

func main() {
	// init config
	config.Init()

	client := objectstorage.NewClient(config.PrivateKey, config.ChainId, config.RpcAddr)

	app := router.SetupRouter(client)

	// Run http service.
	err := app.Run()
	common.HandleErr(err, "Gin run failed.")
}

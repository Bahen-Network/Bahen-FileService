package main

import (
	"file-service/config"
	"file-service/router"
	"file-service/storageclient"
	common "file-service/util"
)

// main.go

func main() {
	// init config
	config.Init()

	client := storageclient.NewClient(config.PrivateKey, config.ChainId, config.RpcAddr)

	app := router.SetupRouter(client)

	// Run cleanup in background after all initializations
	// go storageclient.CleanupOldFiles() // Assuming the function is placed inside storageclient

	// Run http service.
	err := app.Run()
	common.HandleErr(err, "Gin run failed.")
}

package main

import (
	"genesis/pos/presentation/boostrap"
	container "genesis/pos/presentation/container"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(os.Stdout)
	boostrap.UploadEnv()
	container.StartContainer()
	// go workers.StartBroker()
	boostrap.RunServer()
}

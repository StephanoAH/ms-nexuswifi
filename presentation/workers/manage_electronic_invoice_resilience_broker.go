package workers

import (
	"log"
	"time"
)

func StartBroker() {
	for {
		log.Println("__________________________________________________")
		log.Println("Time Instance...", time.Now())
		log.Println("__________________________________________________")
		time.Sleep(10 * time.Minute)
	}
}

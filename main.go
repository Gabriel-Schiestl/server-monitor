package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/Gabriel-Schiestl/server-monitor/internal/domain/model"
	"github.com/Gabriel-Schiestl/server-monitor/internal/infra/service"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading env: %v", err)
	}
	
	logs := make(chan model.Log)
	rmq := model.NewRabbitMQ("email-service")
	defer rmq.Close()

	health := model.Health{
		Servers: []model.Server{
			{
				ID:     1,
				Host:   "http://localhost:3000",
				Status: "ENABLED",
			},
			{
				ID:     2,
				Host:   "https://amazon.com.br",
				Status: "ENABLED",
			},
		},
		Logs: logs,
		LimitTime: 10,
		IntervalTime: 5,
		Wg: sync.WaitGroup{},
		Mu: sync.Mutex{},
	}

	service := service.NewHealthService(&health)
	go service.Check()

	for log := range logs {
		if log.StatusCode != 200 {
			rmq.Publish("nothorizon@hotmail.com", log)
		}
		fmt.Println(log)
	}
}
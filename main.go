package main

import (
	"fmt"
	"sync"

	"github.com/Gabriel-Schiestl/server-monitor/internal/domain/model"
	"github.com/Gabriel-Schiestl/server-monitor/internal/infra/service"
)

func main() {
	logs := make(chan model.Log)

	health := model.Health{
		Servers: []model.Server{
			{
				ID:     1,
				Host:   "http://localhost",
				Port:   8080,
				Status: "ENABLED",
			},
			{
				ID:     2,
				Host:   "http://localhost",
				Port:   8081,
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
		fmt.Println(log)
	}
}
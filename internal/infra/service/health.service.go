package service

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Gabriel-Schiestl/server-monitor/internal/domain/model"
)

type healthService struct {
	health *model.Health
}

func NewHealthService(health *model.Health) *healthService {
	return &healthService{health: health}
}

func (s *healthService) Check() {
	ticker := time.NewTicker(time.Duration(s.health.IntervalTime) * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		s.checkServers()
	}
}

func (s *healthService) checkServers() {
	for _, server := range s.health.Servers {
		url := server.Host + ":" + strconv.Itoa(server.Port)
		date := time.Now().Format("2006-01-02 15:04:05")

		if server.Status != "ENABLED" {
			s.health.Logs <- model.Log{
				ServerID: 	server.ID,
				URL: 		url,
				Status: 	"DISABLED",
				StatusCode: 0,
				Message: 	"Server is disabled",
				ContentLength: 0,
				ElapsedTime: 0,
				DateTime: date,
			}
			continue
		}

		//s.health.Wg.Add(1)
		go func(server model.Server) {
			//defer s.health.Wg.Done()

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				s.health.Logs <- model.Log{
					ServerID: 	server.ID,
					URL: 		url,
					Status: 	"ERROR",
					StatusCode: 0,
					ContentLength: 0,
					ElapsedTime: 0,
					DateTime: date,
				}
				return
			}

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				s.health.Logs <- model.Log{
					ServerID: 	server.ID,
					URL: 		url,
					Status: 	"ERROR",
					StatusCode: 0,
					ContentLength: 0,
					ElapsedTime: 0,
					DateTime: date,
				}
				return
			}

			elapsedTime, err := strconv.ParseFloat(res.Header.Get("X-Runtime"), 64)
			if err != nil {
				elapsedTime = 0
			}

			s.health.Logs <- model.Log{
				ServerID: 	server.ID,
				URL: 		url,
				Status: 	"OK",
				StatusCode: res.StatusCode,
				ContentLength: res.ContentLength,
				ElapsedTime: elapsedTime,
				DateTime: date,
			}
		}(server)
		//s.health.Wg.Wait()
	}
}
package model

type Server struct {
	ID     int    `json:"id"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Status string `json:"status"`
}
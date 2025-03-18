package model

type Log struct {
	ServerID      int
	URL           string
	Status        string
	StatusCode    int
	Message       string
	ContentLength int64
	ElapsedTime   float64
	DateTime      string
}
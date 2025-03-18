package model

import "sync"

type Health struct {
	Servers      []Server `json:"servers"`
	LimitTime    int      `json:"limit_time"`
	IntervalTime int      `json:"interval_time"`
	Wg           sync.WaitGroup
	Mu 		 	 sync.Mutex
	Logs           chan<- Log
}
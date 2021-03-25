package model

import "time"

type Server struct {
	Hostname     string    `json:"hostname"`
	Key          string    `json:"key"`
	IPs          []string  `json:"ips"`
	Dependencies []string  `json:"dependencies"`
	LastUpdate   time.Time `json:"lastUpdate"`
}

type Model interface {
	Add(Server) error
	GetAll() ([]Server, error)
	GetSince(time.Duration) ([]Server, error)
	Clear() error
}

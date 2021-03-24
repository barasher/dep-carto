package model

import "time"

type Server struct {
	Hostname     string
	Key          string
	IPs          []string
	Dependencies []string
	LastUpdate   time.Time
}


type Model interface {
	AddServer(server Server) error
	GetAllServers() ([]Server, error)
}

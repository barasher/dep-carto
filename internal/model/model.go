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
	AddServer(Server) error
	GetAllServers() ([]Server, error)
	GetServersSince(time.Duration) ([]Server, error)
	Clear() error
}

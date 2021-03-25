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
	Add(Server) error
	GetAll() ([]Server, error)
	GetSince(time.Duration) ([]Server, error)
	Clear() error
}

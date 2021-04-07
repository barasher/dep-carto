package model

import (
	"context"
	"time"
)

type Server struct {
	Hostname     string    `json:"hostname"`
	Key          string    `json:"key"`
	IPs          []string  `json:"ips"`
	Dependencies []string  `json:"dependencies"`
	LastUpdate   time.Time `json:"lastUpdate"`
}

type Model interface {
	Add(ctx context.Context, s Server) error
	GetAll(ctx context.Context) ([]Server, error)
	GetAllSince(ctx context.Context, d time.Duration) ([]Server, error)
	Clear(ctx context.Context) error
}

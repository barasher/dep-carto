package model

import (
	"context"
	"time"
)

type Server struct {
	Hostname     string       `json:"hostname"`
	Key          string       `json:"key"`
	IPs          []string     `json:"ips"`
	Dependencies []Dependency `json:"dependencies"`
	LastUpdate   time.Time    `json:"lastUpdate"`
}

type Dependency struct {
	Resource string `json:"resource"`
	Label    string `json:"label"`
}

type Model interface {
	Add(ctx context.Context, s Server) error
	GetAll(ctx context.Context, d *time.Duration) ([]Server, error)
	Clear(ctx context.Context) error
	GetDepending(ctx context.Context, ident string, depth *int, since *time.Duration) ([]Server, error)
	GetDependencies(ctx context.Context, ident string, depth *int, since *time.Duration) ([]Server, error)
}

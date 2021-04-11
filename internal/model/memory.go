package model

import (
	"context"
	"time"
)

type memoryModel struct {
	servers []Server
}

func NewMemoryModel() Model {
	return &memoryModel{
		servers: []Server{},
	}
}

func (m *memoryModel) Add(ctx context.Context, server Server) error {
	for i, s := range m.servers {
		if s.Hostname == server.Hostname && s.Key == server.Key {
			m.servers[i] = server
			return nil
		}
	}
	m.servers = append(m.servers, server)
	return nil
}

func (m *memoryModel) GetAll(ctx context.Context, d *time.Duration) ([]Server, error) {
	if d == nil {
		return m.servers, nil
	}
	limit := time.Now().Add(-*d)
	var s []Server
	for _, curS := range m.servers {
		if !curS.LastUpdate.Before(limit) {
			s = append(s, curS)
		}
	}
	return s, nil
}

func (m *memoryModel) Clear(ctx context.Context) error {
	m.servers = []Server{}
	return nil
}

func (m *memoryModel) GetDepending(ctx context.Context, ident string, depth *int, since *time.Duration) ([]Server, error) {
	return []Server{}, nil
}

func (m *memoryModel) GetDependencies(ctx context.Context, ident string, depth *int, since *time.Duration) ([]Server, error) {
	return []Server{}, nil
}

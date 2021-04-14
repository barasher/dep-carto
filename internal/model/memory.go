package model

import (
	"context"
	"fmt"
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

func (m *memoryModel) getByIdentifier(ident string) []Server {
	var s []Server
	for _, c := range m.servers {
		if c.IsIdentifiedBy(ident) {
			s = append(s, c)
		}
	}
	return s
}

func (m *memoryModel) GetDepending(ctx context.Context, ident string, depth int, since time.Duration) ([]Server, error) {
	return []Server{}, fmt.Errorf("not yet implemented")
}

func (m *memoryModel) getDependencies(idents []string) []string {
	dedup := make(map[string]bool)
	for _, curId := range idents {
		for _, curS := range m.getByIdentifier(curId) {
			for _, curD := range curS.Dependencies {
				dedup[curD.Resource] = true
			}
		}
	}
	deps := make([]string, len(dedup))
	i := 0
	for k, _ := range dedup {
		deps[i] = k
		i++
	}
	return deps
}

func (m *memoryModel) GetDependencies(ctx context.Context, ident string, depth int, since time.Duration) ([]Server, error) {
	collected := make(map[string]bool)
	collected[ident] = true
	toCrawl := []string{ident}
	for i := 0; i < depth; i++ {
		toContinue := []string{}
		for _, crawledDep := range m.getDependencies(toCrawl) {
			if _, found := collected[crawledDep]; !found {
				toContinue = append(toContinue, crawledDep)
				collected[crawledDep] = true
			}
		}
		if len(toContinue) == 0 {
			break
		}
		toCrawl = toContinue
	}
	var deps []Server
	hostMap := make(map[string]bool)
	limit := time.Now().Add(-since)
	for c, _ := range collected {
		cur := m.getByIdentifier(c)
		if len(cur) > 0 {
			if _, found := hostMap[cur[0].Hostname]; !found {
				for _, curS := range cur {
					hostMap[cur[0].Hostname] = true
					if !curS.LastUpdate.Before(limit) {
						deps = append(deps, curS)
					}
				}
			}
		}
	}
	return deps, nil
}

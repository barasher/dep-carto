package model

import "time"

type memoryModel struct {
	servers []Server
}

func NewMemoryModel() Model {
	return &memoryModel{}
}

func (m *memoryModel) Add(server Server) error {
	for i, s := range m.servers {
		if s.Hostname == server.Hostname && s.Key == server.Key {
			m.servers[i] = server
			return nil
		}
	}
	m.servers = append(m.servers, server)
	return nil
}

func (m *memoryModel) GetAll() ([]Server, error) {
	return m.servers, nil
}

func (m *memoryModel) GetSince(d time.Duration) ([]Server, error) {
	limit := time.Now().Add(-d)
	var s []Server
	for _, curS := range m.servers {
		if ! curS.LastUpdate.Before(limit) {
			s = append(s, curS)
		}
	}
	return s, nil
}

func (m *memoryModel) Clear() error {
	m.servers = []Server{}
	return nil
}





package model

type memoryModel struct {
	servers []Server
}

func NewMemoryModel() Model {
	return &memoryModel{}
}

func (m *memoryModel) AddServer(server Server) error {
	for i, s := range m.servers {
		if s.Hostname == server.Hostname && s.Key == server.Key {
			m.servers[i] = server
			return nil
		}
	}
	m.servers = append(m.servers, server)
	return nil
}

func (m *memoryModel) GetAllServers() ([]Server, error) {
	return m.servers, nil
}





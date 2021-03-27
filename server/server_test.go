package server

import (
	"context"
	"github.com/barasher/dep-carto/internal/model"
	"time"
)

type modelMock struct {
	add struct {
		inServer model.Server
		outErr   error
	}
	getAll struct {
		servers []model.Server
		err     error
	}
	getSince struct {
		inDuration time.Duration
		outServers []model.Server
		outErr     error
	}
	clear struct {
		err error
	}
}

func (m *modelMock) MockAdd(err error) *modelMock {
	m.add.outErr = err
	return m
}

func (m *modelMock) Add(ctx context.Context, s model.Server) error {
	m.add.inServer = s
	return m.add.outErr
}

func (m *modelMock) MockGetAll(s []model.Server, err error) *modelMock {
	m.getAll.servers = s
	m.getAll.err = err
	return m
}

func (m *modelMock) GetAll(ctx context.Context) ([]model.Server, error) {
	return m.getAll.servers, m.getAll.err
}

func (m *modelMock) MockGetSince(s []model.Server, err error) *modelMock {
	m.getSince.outServers = s
	m.getSince.outErr = err
	return m
}

func (m *modelMock) GetSince(ctx context.Context, d time.Duration) ([]model.Server, error) {
	m.getSince.inDuration = d
	return m.getSince.outServers, m.getSince.outErr
}

func (m *modelMock) MockClear(err error) *modelMock {
	m.clear.err = err
	return m
}

func (m *modelMock) Clear(ctx context.Context) error {
	return m.clear.err
}

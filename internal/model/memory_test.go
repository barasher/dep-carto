package model

import "testing"

func TestMemoryModel(t *testing.T) {
	m := NewMemoryModel()
	testModelWorkflow(t, m)
}

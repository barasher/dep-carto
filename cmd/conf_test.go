package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConf_Nominal(t *testing.T) {
	c, err := loadConf("../testdata/conf/server_mem_nominal.json")
	assert.Nil(t, err)
	assert.Equal(t, uint(123), c.Server.Port)
	assert.Equal(t, "memory", c.Server.StorageType)
}

func TestLoadConf_NonExistingFile(t *testing.T) {
	_, err := loadConf("nonExistingFile")
	assert.NotNil(t, err)
}

func TestPort(t *testing.T) {
	c := Conf{Server: ServerConf{Port: 123}}
	assert.Equal(t, uint(123), c.Server.port())
	c = Conf{}
	assert.Equal(t, defaultPort, c.Server.port())
}

func TestStorageType(t *testing.T) {
	c := Conf{Server: ServerConf{StorageType: "blabla"}}
	assert.Equal(t, "blabla", c.Server.storageType())
	c = Conf{}
	assert.Equal(t, defaultStorageType, c.Server.storageType())

}

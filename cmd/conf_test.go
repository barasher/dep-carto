package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConf(t *testing.T) {
	c, err := loadConf("../testdata/conf/server_mem_nominal.json")
	assert.Nil(t, err)
	assert.Equal(t, uint(123), c.Server.Port)
	assert.Equal(t, "memory", c.Server.StorageType)
	assert.Equal(t, "u", c.Crawler.ServerURL)
	assert.ElementsMatch(t, []string{"s1", "s2"}, c.Crawler.Suffixes)
	assert.ElementsMatch(t, []string{"i1", "i2"}, c.Crawler.Inputs)

	_, err = loadConf("nonExistingFile")
	assert.NotNil(t, err)

	_, err = loadConf("../testdata/conf/unparsable.json")
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

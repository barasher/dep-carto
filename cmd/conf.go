package cmd

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	defaultPort        = uint(8080)
	defaultStorageType = memoryStorage

	memoryStorage = "memory"
)

type Conf struct {
	Server ServerConf `json:"server"`
}

type ServerConf struct {
	Port        uint   `json:"port"`
	StorageType string `json:"storage"`
}

func (sc ServerConf) port() uint {
	if sc.Port == 0 {
		return defaultPort
	}
	return sc.Port
}

func (sc ServerConf) storageType() string {
	if sc.StorageType == "" {
		return defaultStorageType
	}
	return sc.StorageType
}

func loadConf(f string) (Conf, error) {
	c := Conf{}
	r, err := os.Open(f)
	if err != nil {
		return c, fmt.Errorf("error while opening configuration file %v: %w", f, err)
	}
	defer r.Close()
	err = json.NewDecoder(r).Decode(&c)
	if err != nil {
		return c, fmt.Errorf("error while unmarshaling configuration file %v: %w", f, err)
	}
	return c, nil
}

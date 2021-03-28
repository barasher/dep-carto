package cmd

import (
	"fmt"
	"github.com/barasher/dep-carto/internal/model"
	"github.com/barasher/dep-carto/server"
	"github.com/spf13/cobra"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "dep-carto : server",
		RunE:  execServer,
	}
)

func init() {
	serverCmd.Flags().StringVarP(&confFile, "conf", "c", "", "dep-carto configuration file")
	serverCmd.MarkFlagRequired("conf")
	RootCmd.AddCommand(serverCmd)
}

func execServer(cmd *cobra.Command, args []string) error {
	c, err := loadConf(confFile)
	if err != nil {
		return err
	}
	m, err := getModel(c)
	if err != nil {
		return err
	}
	s, err := server.NewServer(m, c.Server.port())
	if err != nil {
		return err
	}
	return s.Run()
}

func getModel(c Conf) (model.Model, error) {
	switch c.Server.storageType() {
	case memoryStorage:
		return model.NewMemoryModel(), nil
	default:
		return nil, fmt.Errorf("unsupported storage type (%v)", c.Server.storageType())
	}
}

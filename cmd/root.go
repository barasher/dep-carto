package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

const (
	retOk          int = 0
	retConfFailure int = 1
	retExecFailure int = 2
)

var (
	RootCmd = &cobra.Command{
		Use:   "dep-carto",
		Short: "Dependency cartographer",
	}
	confFile string
)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Error().Msgf("%v", err)
		os.Exit(1)
	}
}
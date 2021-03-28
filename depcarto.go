package main

import (
	"github.com/barasher/dep-carto/cmd"
	"github.com/rs/zerolog/log"
	"os"
)

const (
	retOk int = 0
	retKo int = 1
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Error().Msgf("%v", err)
		os.Exit(retKo)
	}
	os.Exit(retOk)
}

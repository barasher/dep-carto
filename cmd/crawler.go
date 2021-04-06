package cmd

import (
	"fmt"

	"github.com/barasher/dep-carto/internal/crawler"
	"github.com/barasher/dep-carto/internal/parser"

	"github.com/spf13/cobra"
)

var (
	crawlerCmd = &cobra.Command{
		Use:   "crawler",
		Short: "dep-carto : crawler",
		RunE:  execCrawler,
	}
)

func init() {
	crawlerCmd.Flags().StringVarP(&confFile, "conf", "c", "", "dep-carto configuration file")
	crawlerCmd.MarkFlagRequired("conf")
	crawlerCmd.Flags().StringArrayVarP(&inputsParam, "input", "i", []string{}, "input files/folders to crawl")
	RootCmd.AddCommand(crawlerCmd)
}

func execCrawler(cmd *cobra.Command, args []string) error {
	c, err := loadConf(confFile)
	if err != nil {
		return err
	}

	var reOpts []parser.RefExtractorOption
	for _, cur := range c.Crawler.Suffixes {
		reOpts = append(reOpts, parser.WithSuffix(cur))
	}
	re, err := parser.NewRefExtractor(reOpts...)
	if err != nil {
		return fmt.Errorf("error while initializing RefExtractor: %w", err)
	}

	cr, err := crawler.NewCrawler(c.Crawler.ServerURL, *re)
	if err != nil {
		return fmt.Errorf("error while initializing crawler: %w", err)
	}

	in := append(c.Crawler.Inputs, inputsParam...)
	if err := cr.Crawl(in); err != nil {
		return fmt.Errorf("error while crawling files: %w", err)
	}
	return nil
}

package cmd

import (
	"log"
	"github.com/spf13/cobra"
	"github.com/Lavos/gbvideo"
	"github.com/Lavos/gbvideo/storers"
	"github.com/kelseyhightower/envconfig"
)

var (
	s gbvideo.Storer
	q gbvideo.Queuer

	c Configuration
)

type Configuration struct {
	DatabaseLocation string
	DownloadLocation string
	APIKey string
}

var RootCmd = &cobra.Command{
	Use:   "gbv",
	Short: "Sync, List, Queue and Download Giantbomb videos.",
}

func init(){
	envconfig.Process("gbvideo", &c)

	var err error
	sl, err := storers.NewSQLite(c.DatabaseLocation)
	s = sl
	q = sl

	if err != nil {
		log.Fatal(err)
	}
}

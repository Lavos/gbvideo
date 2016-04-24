package cmd

import (
	"log"
	"time"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/Lavos/gbvideo/giantbomb"
	"github.com/daviddengcn/go-colortext"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Syncs the local database with the Giantbomb API.",
	RunE: func(cmd *cobra.Command, args []string) error {
		gb := giantbomb.New(giantbomb.ProductionAPILocation, c.APIKey)

		count, err := s.GetCount()

		if err != nil {
			return err
		}

		vr, err := gb.GetVideos(0, 100, giantbomb.Id, giantbomb.DirectionDesc)

		if err != nil {
			return err
		}

		for _, v := range vr.Results {
			err = s.InsertVideo(v)

			if err != nil {
				log.Printf("INSERT ERROR: %#v", err)
			}
		}

		fmt.Printf("Starting Sync operation...\n")
		fmt.Printf("Videos in database: %d\n", count)
		fmt.Printf("Total Giant Bomb videos: %d\n", vr.TotalResults)
		fmt.Printf("Missing %d videos in database.\n", vr.TotalResults - count)

		if count < vr.TotalResults {
			ct.ChangeColor(ct.Yellow, true, ct.Black, false)
			fmt.Printf("Backscanning for missing videos...\n")
			ct.ResetColor()
		}

		var offset int64
		for count < vr.TotalResults {
			time.Sleep(1500 * time.Millisecond)

			vr, err = gb.GetVideos(offset, 100, giantbomb.Id, giantbomb.DirectionDesc)

			fmt.Printf("Got %d videos.\n", vr.PageResults)

			if err != nil {
				return err
			}

			for _, v := range vr.Results {
				err = s.InsertVideo(v)

				if err != nil {
					log.Printf("INSERT ERROR: %#v", err)
				}
			}

			count, err = s.GetCount()

			if err != nil {
				return err
			}

			offset += int64(len(vr.Results))
		}

		return nil
	},
}

func init(){
	RootCmd.AddCommand(syncCmd)
}

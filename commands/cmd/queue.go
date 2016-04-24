package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/daviddengcn/go-colortext"
)

var queueCmd = &cobra.Command{
	Use:   "queue",
	Short: "Lists videos marked for download.",
	RunE: func(cmd *cobra.Command, args []string) error {
		videos, err := q.GetQueuedVideos()

		if err != nil {
			return err
		}

		ct.ChangeColor(ct.Magenta, true, ct.Black, false)
		fmt.Printf("Current queued videos: %d\n", len(videos))
		ct.ResetColor()

		for _, v := range videos {
			printVideo(v)
		}

		return nil
	},
}

func init(){
	RootCmd.AddCommand(queueCmd)
}

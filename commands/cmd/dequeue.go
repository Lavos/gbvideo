package cmd

import (
	"fmt"
	"strconv"
	"github.com/spf13/cobra"
	"github.com/Lavos/gbvideo"
)

var dequeueCmd = &cobra.Command{
	Use:   "dequeue",
	Short: "Unmarks a video for download.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var vd *gbvideo.VideoDownload
		var id_num int64
		var err error

		for _, a := range args {
			id_num, err = strconv.ParseInt(a, 10, 64)

			if err != nil {
				return err
			}

			vd, err = s.GetVideo(id_num)

			if err != nil {
				return err
			}

			err = q.Dequeue(vd)

			if err != nil {
				return err
			}

			fmt.Printf("Dequeued %s successfully.\n", a)
		}

		return nil
	},
}

func init(){
	RootCmd.AddCommand(dequeueCmd)
}

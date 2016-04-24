package cmd

import (
	"fmt"
	"strconv"
	"github.com/spf13/cobra"
	"github.com/Lavos/gbvideo"
)

var enqueueCmd = &cobra.Command{
	Use:   "enqueue",
	Short: "Marks a video for download.",
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

			err = q.Enqueue(vd)

			if err != nil {
				return err
			}

			fmt.Printf("Enqueued %s successfully.\n", a)
		}

		return nil
	},
}

func init(){
	RootCmd.AddCommand(enqueueCmd)
}

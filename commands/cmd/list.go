package cmd

import (
	"fmt"
	"strconv"
	"github.com/spf13/cobra"
	"github.com/Lavos/gbvideo/storers"
)

var (
	offset int64
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists the specified number of videos in the database.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("No number of videos specified.")
		}

		num, err := strconv.ParseInt(args[0], 10, 64)

		if err != nil {
			return err
		}

		videos, err := s.GetVideos(offset, num, storers.Field_PublishDate, false)

		if err != nil {
			return err
		}

		for _, v := range videos {
			printVideo(v)
		}

		return nil
	},
}

func init(){
	listCmd.Flags().Int64VarP(&offset, "offset", "o", 0, "offset to start list iteration")
	RootCmd.AddCommand(listCmd)
}

package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"log"
	"os/exec"
	"github.com/spf13/cobra"
	"github.com/Lavos/gbvideo"
	"github.com/daviddengcn/go-colortext"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads queued videos.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		videos, err := q.GetQueuedVideos()

		if err != nil {
			return err
		}

		var pr *gbvideo.ProgressReader
		var video_url *url.URL
		var values url.Values
		var req *http.Request
		var resp *http.Response
		done := make(chan bool)

		for _, video := range videos {
			// create URL
			video_url, err = url.Parse(video.HighURL)

			if err != nil {
				return err
			}

			values = video_url.Query()
			values.Set("api_key", c.APIKey)

			video_url.RawQuery = values.Encode()

			// open http request
			req, err = http.NewRequest("GET", video_url.String(), nil)

			if err != nil {
				return err
			}

			resp, err = http.DefaultClient.Do(req)

			if err != nil {
				log.Printf("HTTP error when downloading `%s`: %s", video.HighURL, err)
				continue
			}

			if resp.StatusCode != 200 {
				log.Printf("non-200 when downloading `%s`: %s", video.HighURL, resp.Status)
				continue
			}

			pr = gbvideo.NewProgressReader(resp.Body)

			go printBytes(pr.BytesRead, resp.ContentLength, done)
			defer resp.Body.Close()

			args := []string{
				"-y",
				"-i", "pipe:",
				"-c:v", "copy",
				"-c:a", "copy",
				"-f", "hls",
				"-threads", "0",
				"-hls_flags", "single_file",
				"-hls_list_size", "0",
				"-bsf:v", "h264_mp4toannexb",
				fmt.Sprintf("%s/%s.m3u8", c.DownloadLocation, video.FileName),
			}

			cmd := exec.Command("ffmpeg", args...)
			cmd.Stdin = pr

			ct.ChangeColor(ct.Green, false, ct.Black, false)
			fmt.Printf("Downloading: ")
			ct.ResetColor()
			fmt.Printf("%s\n", video.FileName)

			err = cmd.Start()

			if err != nil {
				log.Printf("Could not exec ffmpeg: %s", err)
				continue
			}

			<-done
			cmd.Wait()

			err = q.MarkDownloaded(video)

			if err != nil {
				log.Printf("Could not mark `%d` as downloaded: %s", video.ID, err)
				continue
			}
		}

		return nil
	},
}

func init(){
	RootCmd.AddCommand(downloadCmd)
}

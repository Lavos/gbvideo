package main

import (
	"net/http"
	"github.com/Lavos/gbvideo"
	"os/exec"
	"log"

	"github.com/cheggaaa/pb"
)

var (
	args = []string{
		"-y",
		"-i", "pipe:",
		"-c:v", "copy",
		"-c:a", "copy",
		"-f", "hls",
		"-threads", "0",
		"-hls_flags", "single_file",
		"-hls_list_size", "0",
		"-bsf:v", "h264_mp4toannexb",
		"out.m3u8",
	}
)

func printBytes (bytesRead chan int, total int64, done chan bool) {
	var r int
	var n int64
	var ok bool

	bar := pb.New64(total)
	bar.ShowSpeed = true
	bar.ShowFinalTime = true
	bar.SetUnits(pb.U_BYTES)

	bar.Start()

	for {
		r, ok = <-bytesRead

		n += int64(r)
		bar.Set64(n)

		if !ok {
			break
		}
	}

	bar.Finish()
	done <- true
}

func main () {
	req, err := http.NewRequest("GET", "http://127.0.0.1:9016/video.mp4", nil)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	pr := gbvideo.NewProgressReader(resp.Body)
	done := make(chan bool)

	go printBytes(pr.BytesRead, resp.ContentLength, done)
	defer resp.Body.Close()

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdin = pr
	cmd.Start()

	<-done
	cmd.Wait()
}

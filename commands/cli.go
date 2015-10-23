package main

import (
	"fmt"
	"log"
	"flag"
	"strings"
	"strconv"
	"time"
	"os"
	"io"
	"net/http"
	"github.com/Lavos/gbvideo"
	"github.com/Lavos/gbvideo/giantbomb"
	"github.com/Lavos/gbvideo/storers"
	"github.com/daviddengcn/go-colortext"
	"github.com/cheggaaa/pb"
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

func sync () {
	gb := giantbomb.New(giantbomb.ProductionAPILocation, c.APIKey)

	count, err := s.GetCount()

	if err != nil {
		log.Fatal(err)
	}

	vr, err := gb.GetVideos(0, 1, giantbomb.Id, giantbomb.DirectionDesc)

	if err != nil {
		log.Fatal(err)
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
		if offset > 0 {
			fmt.Printf("Sleeping 12secs (5calls/min)...\n")
			time.Sleep(12 * time.Second)
		}

		vr, err = gb.GetVideos(offset, 100, giantbomb.Id, giantbomb.DirectionDesc)

		fmt.Printf("Got %d videos.\n", vr.PageResults)

		if err != nil {
			log.Fatal(err)
		}

		for _, v := range vr.Results {
			err = s.InsertVideo(v)

			if err != nil {
				log.Printf("INSERT ERROR: %#v", err)
			}
		}

		count, err = s.GetCount()

		if err != nil {
			log.Fatal(err)
		}

		offset += int64(len(vr.Results))
	}
}

func top (num int64) {
	fmt.Printf("Top %d videos by publish date:\n", num)
	videos, err := s.GetVideos(num, storers.Field_PublishDate, false)

	if err != nil {
		log.Fatal(err)
	}

	for _, v := range videos {
		ct.ChangeColor(ct.Yellow, true, ct.Black, false)
		fmt.Printf("%d: ", v.ID)
		ct.ChangeColor(ct.White, true, ct.Black, false)
		fmt.Printf("%s", v.Name)

		ct.ChangeColor(ct.Cyan, true, ct.Black, false)
		fmt.Printf(" [%s] ", v.VideoType)

		if v.DownloadDate != nil {
			ct.ChangeColor(ct.Green, true, ct.Black, false)
			fmt.Printf("D")
		}

		if v.Queued {
			ct.ChangeColor(ct.Blue, true, ct.Black, false)
			fmt.Printf("Q")
		}

		ct.ChangeColor(ct.White, false, ct.Black, false)
		fmt.Printf("\n\t%s\n", v.Deck)
	}
}

func queue (nums ...int64) {
	var vd *gbvideo.VideoDownload
	var err error

	for _, n := range nums {
		vd, err = s.GetVideo(n)

		if err != nil {
			log.Fatal(err)
		}

		err = q.Queue(vd)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Queued %d successfully.\n", n)
	}
}

func unqueue (nums ...int64) {
	var vd *gbvideo.VideoDownload
	var err error

	for _, n := range nums {
		vd, err = s.GetVideo(n)

		if err != nil {
			log.Fatal(err)
		}

		err = q.UnQueue(vd)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("UnQueued %d successfully.\n", n)
	}
}

func download() {
	var err error

	videos, err := q.GetQueuedVideos()

	if err != nil {
		log.Fatal(err)
	}

	var pr *gbvideo.ProgressReader
	var req *http.Request
	var resp *http.Response
	var file *os.File
	done := make(chan bool)

	for _, video := range videos {
		// open file
		file, err = os.Create(fmt.Sprintf("%s/%s", c.DownloadLocation, video.FileName))

		if err != nil {
			log.Fatal(err)
		}

		// open http request
		req, err = http.NewRequest("GET", video.HighURL, nil)

		if err != nil {
			log.Fatal(err)
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
		defer file.Close()

		ct.ChangeColor(ct.Green, false, ct.Black, false)
		fmt.Printf("Downloading: ")
		ct.ResetColor()
		fmt.Printf("%s\n", video.FileName)

		_, err = io.Copy(file, pr)

		if err != nil {
			log.Fatalf("COPY ERROR: %s", err)
		}

		<-done

		err = q.MarkDownloaded(video)

		if err != nil {
			log.Printf("Could not mark `%d` as downloaded: %s", video.ID, err)
			continue
		}
	}
}

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

func splitNumbers (str_list string) []int64 {
	str_slice := strings.Split(str_list, ",")
	num_list := make([]int64, len(str_slice))
	var counter int

	for _, str := range str_slice {
		m, err := strconv.ParseInt(str, 10, 64)

		if err != nil {
			continue
		}

		num_list[counter] = m
		counter++
	}

	return num_list[:counter]
}

func main () {
	flag.Parse()
	envconfig.Process("gbvideo", &c)

	var err error
	sl, err := storers.NewSQLite(c.DownloadLocation)
	s = sl
	q = sl

	if err != nil {
		log.Fatal(err)
	}

	command := flag.Arg(0)

	switch command {
	case "sync":
		sync()

	case "top":
		var i int64
		count := flag.Arg(1)

		if count == "" {
			i = 10
		} else {
			i, err = strconv.ParseInt(count, 10, 64)

			if err != nil {
				i = 10
			}
		}

		top(i)

	case "queue":
		nums := splitNumbers(flag.Arg(1))
		queue(nums...)

	case "unqueue":
		nums := splitNumbers(flag.Arg(1))
		unqueue(nums...)

	case "download":
		download()
	}
}

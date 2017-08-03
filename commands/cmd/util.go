package cmd

import (
	"fmt"
	"strings"
	"strconv"
	"github.com/Lavos/gbvideo"
	"github.com/cheggaaa/pb"
	"github.com/daviddengcn/go-colortext"
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

func printVideo (v *gbvideo.VideoDownload) {
	ct.ChangeColor(ct.Yellow, true, ct.Black, false)
	fmt.Printf("%d: ", v.ID)
	ct.ChangeColor(ct.White, true, ct.Black, false)
	fmt.Printf("%s", v.Name)

	ct.ChangeColor(ct.Cyan, true, ct.Black, false)
	fmt.Printf(" [%s] ", v.VideoType)

	ct.ChangeColor(ct.Red, false, ct.Black, false)
	fmt.Printf("(%s) ", v.FileName)

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

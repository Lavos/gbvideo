package gbvideo

import (
	"time"

	"github.com/Lavos/gbvideo/giantbomb"
)

type VideoDownload struct {
	giantbomb.Video

	DownloadDate *time.Time
	Queued bool
}

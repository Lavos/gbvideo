package gbvideo

import (
	"github.com/Lavos/gbvideo/giantbomb"
)

type Storer interface {
	GetCount() (int64, error)
	GetVideo(id int64) (*VideoDownload, error)
	GetVideos(offset, limit int64, sort_field int, direction_asce bool) ([]*VideoDownload, error)
	InsertVideo(*giantbomb.Video) error
}

type Queuer interface {
	Enqueue (*VideoDownload) error
	Dequeue (*VideoDownload) error
	MarkDownloaded (*VideoDownload) error
	GetQueuedVideos() ([]*VideoDownload, error)
}

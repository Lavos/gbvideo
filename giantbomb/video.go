package giantbomb

import (
	"time"
)

const (
	PublishDate = "publish_date"
	Id = "id"
)

type VideoResponse struct {
	Response

	Results []*Video
}

type JSONDate struct {
	time.Time
}

type Video struct {
	ID int64 `json:"id"`

	APIDetailURL string `json:"api_detail_url"`
	SiteDetailURL string `json:"site_detail_url"`

	Deck string `json:"deck"`
	HighURL string `json:"high_url"`

	PublishDate *JSONDate `json:"publish_date"`

	Name string `json:"name"`
	Length int64 `json:"length_seconds"`
	FileName string `json:"url"`
	VideoType string `json:"video_type"`
}

// 2015-10-17 06:00:00
// 2006-01-02 15:04:05
func (j *JSONDate) UnmarshalJSON(p []byte) error {
	t, err := time.Parse(`"2006-01-02 15:04:05"`, string(p))

	if err != nil {
		return err
	}

	j.Time = t
	return nil
}

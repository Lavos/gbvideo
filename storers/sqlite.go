package storers

import (
	"fmt"
	"time"
	"github.com/Lavos/gbvideo"
	"github.com/Lavos/gbvideo/giantbomb"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	Field_ID = iota
	Field_APIDetailURL
	Field_SiteDetailURL
	Field_Deck
	Field_HighURL
	Field_PublishDate
	Field_Name
	Field_Length
)

var (
	field_map = map[int]string{
		0: "id",
		1: "api_detail_url",
		2: "site_detail_url",
		3: "deck",
		4: "high_url",
		5: "publish_date",
		6: "name",
		7: "length",
	}
)

type SQLite struct {
	db *sql.DB
}

func NewSQLite (location string) (*SQLite, error) {
	db, err := sql.Open("sqlite3", location)

	if err != nil {
		return nil, err
	}

	return &SQLite{
		db: db,
	}, nil
}

func (s *SQLite) getVideosFromRows(rows *sql.Rows) ([]*gbvideo.VideoDownload, error) {
	videos := make([]*gbvideo.VideoDownload, 0)

	var id, publish_date, length int64
	var api_detail_url, site_detail_url, deck, high_url, name, filename, video_type string
	var video *gbvideo.VideoDownload

	var queue_id, download_ts sql.NullInt64

	var t, t2 time.Time
	var d *time.Time

	var err error

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &api_detail_url, &site_detail_url, &deck, &high_url, &publish_date, &download_ts, &name, &length, &filename, &video_type, &queue_id)

		if err != nil {
			continue
		}

		t = time.Unix(publish_date, 0)

		if download_ts.Valid {
			t2 = time.Unix(download_ts.Int64, 0)
			d = &t2
		} else {
			d = nil
		}

		video = &gbvideo.VideoDownload{
			Video: giantbomb.Video{
				ID: id,
				APIDetailURL: api_detail_url,
				SiteDetailURL: site_detail_url,
				Deck: deck,
				HighURL: high_url,
				PublishDate: &giantbomb.JSONDate{t},
				Name: name,
				Length: length,
				FileName: filename,
				VideoType: video_type,
			},

			DownloadDate: d,
			Queued: queue_id.Valid,
		}

		videos = append(videos, video)
	}

	return videos, nil
}

func (s *SQLite) GetCount () (int64, error) {
	var count int64
	err := s.db.QueryRow("SELECT COUNT(id) FROM videos").Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *SQLite) InsertVideo (video *giantbomb.Video) error {
	_, err := s.db.Exec(`INSERT OR IGNORE INTO videos (id, api_detail_url, site_detail_url, deck, high_url, publish_date, name, length, filename, video_type) VALUES (?, ?, ?, ?, ?,  ?, ?, ?, ?, ?)`,
		video.ID, video.APIDetailURL, video.SiteDetailURL, video.Deck, video.HighURL, video.PublishDate.Unix(), video.Name, video.Length, video.FileName, video.VideoType)

	return err
}

func (s *SQLite) GetVideo(id int64) (*gbvideo.VideoDownload, error) {
	rows, err := s.db.Query("SELECT * FROM videos LEFT JOIN queue ON videos.id = queue.video_id WHERE videos.id = ? LIMIT 1", id)

	if err != nil {
		return nil, err
	}

	videos, err := s.getVideosFromRows(rows)

	if err != nil {
		return nil, err
	}

	if len(videos) != 1 {
		return nil, fmt.Errorf("No video found for that ID.")
	}

	return videos[0], nil
}

func (s *SQLite) GetVideos(limit int64, sort_field int, direction_asce bool) ([]*gbvideo.VideoDownload, error) {
	key, ok := field_map[sort_field]

	if !ok {
		return nil, fmt.Errorf("Invalid sort field key.")
	}

	var direction string

	if direction_asce {
		direction = "ASC"
	} else {
		direction = "DESC"
	}

	query := fmt.Sprintf("SELECT * FROM videos LEFT JOIN queue ON videos.id = queue.video_id ORDER BY %s %s LIMIT %d", key, direction, limit)
	rows, err := s.db.Query(query)

	switch {
	case err == sql.ErrNoRows:
		return make([]*gbvideo.VideoDownload, 0), nil

	case err != nil:
		return nil, err
	}

	return s.getVideosFromRows(rows)
}

func (s *SQLite) Enqueue(video *gbvideo.VideoDownload) error {
	_, err := s.db.Exec(`INSERT OR IGNORE INTO queue (video_id) VALUES (?)`, video.ID)

	return err
}

func (s *SQLite) Dequeue(video *gbvideo.VideoDownload) error {
	_, err := s.db.Exec(`DELETE FROM queue WHERE video_id = ?`, video.ID)

	return err
}

func (s *SQLite) MarkDownloaded(video *gbvideo.VideoDownload) error {
	_, err := s.db.Exec(`DELETE FROM queue WHERE video_id = ?`, video.ID)

	if err != nil {
		return err
	}

	now_ts := time.Now().Unix()
	_, err = s.db.Exec(`UPDATE videos SET download_date = ? WHERE id = ?`, now_ts, video.ID)

	return err
}

func (s *SQLite) GetQueuedVideos() ([]*gbvideo.VideoDownload, error) {
	query := fmt.Sprintf("SELECT * FROM videos LEFT JOIN queue ON videos.id = queue.video_id WHERE videos.id = queue.video_id")
	rows, err := s.db.Query(query)

	switch {
	case err == sql.ErrNoRows:
		return make([]*gbvideo.VideoDownload, 0), nil

	case err != nil:
		return nil, err
	}

	return s.getVideosFromRows(rows)
}

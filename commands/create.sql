CREATE TABLE videos (
	id INTEGER PRIMARY KEY,
	api_detail_url TEXT NOT NULL,
	site_detail_url TEXT NOT NULL,
	deck TEXT NOT NULL,
	high_url TEXT NOT NULL,
	publish_date INTEGER NOT NULL,
	download_date INTEGER,
	name TEXT NOT NULL,
	length INTEGER NOT NULL,
	filename TEXT NOT NULL,
	video_type TEXT NOT NULL
);

CREATE TABLE queue (
	video_id INTEGER UNIQUE NOT NULL
);

CREATE INDEX id ON videos (id);
CREATE INDEX publish_date ON videos (publish_date);
CREATE INDEX video_id ON queue (video_id);

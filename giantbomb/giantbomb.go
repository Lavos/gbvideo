package giantbomb

import (
	"fmt"
	"net/url"
	"net/http"
	"encoding/json"
	"strconv"
)

const (
	ProductionAPILocation string = "http://www.giantbomb.com/api"

	DirectionAsce = "asc"
	DirectionDesc = "desc"
)

var (
	StatusCodeMap = map[int64]string {
		1: "OK",
		100: "Invalid API Key",
		101: "Object Not Found",
		102: "Error in URL Format",
		103: "'JSONP' format requires a `json_callback` argument",
		104: "Filter Error",
		105: "Subscriber only video is for subscribers only.",
	}
)

type GiantBomb struct {
	apiLocation string
	apiKey string
}

func New (api_location, api_key string) *GiantBomb {
	return &GiantBomb{
		apiLocation: api_location,
		apiKey: api_key,
	}
}

func (g *GiantBomb) GetVideos (offset, limit int64, sort_field, direction string) (*VideoResponse, error) {
	u, err := url.Parse(fmt.Sprintf("%s/videos", g.apiLocation))

	if err != nil {
		return nil, err
	}

	v := u.Query()
	v.Add("format", "json")
	v.Add("api_key", g.apiKey)
	v.Add("offset", strconv.FormatInt(offset, 10))
	v.Add("limit", strconv.FormatInt(limit, 10))
	v.Add("sort", fmt.Sprintf("%s:%s", sort_field, direction))
	u.RawQuery = v.Encode()

	request, err := http.NewRequest("GET", u.String(), nil)

	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Got a non-200 error: %s", response.Status)
	}

	var vr VideoResponse
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&vr)

	if err != nil {
		return nil, err
	}

	if vr.StatusCode != 1 {
		status_err, ok := StatusCodeMap[vr.StatusCode]

		if ok {
			return nil, fmt.Errorf(status_err)
		}

		return nil, fmt.Errorf("Unknown error returned. Status code: %d", vr.StatusCode)
	}

	return &vr, nil
}

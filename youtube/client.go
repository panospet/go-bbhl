package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	youtubeApiUrl = "https://www.googleapis.com/youtube/v3/search"
)

type Client struct {
	apiKey string
	http   *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		http:   http.DefaultClient,
	}
}

type VideoInfo struct {
	Id          string
	Title       string
	PublishedAt time.Time
}

type YoutubeResponse struct {
	Kind          string `json:"kind"`
	Etag          string `json:"etag"`
	NextPageToken string `json:"nextPageToken"`
	RegionCode    string `json:"regionCode"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []Item `json:"items"`
}

type Item struct {
	Kind string `json:"kind"`
	Etag string `json:"etag"`
	Id   struct {
		Kind    string `json:"kind"`
		VideoId string `json:"videoId"`
	} `json:"id"`
	Snippet struct {
		PublishedAt time.Time `json:"publishedAt"`
		ChannelId   string    `json:"channelId"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Thumbnails  struct {
			Default struct {
				Url    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"default"`
			Medium struct {
				Url    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"medium"`
			High struct {
				Url    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"high"`
		} `json:"thumbnails"`
		ChannelTitle         string    `json:"channelTitle"`
		LiveBroadcastContent string    `json:"liveBroadcastContent"`
		PublishTime          time.Time `json:"publishTime"`
	} `json:"snippet"`
}

func (c *Client) GetChannelVideos(
	ctx context.Context,
	channelId string,
	maxResults int,
) ([]VideoInfo, error) {
	url := fmt.Sprintf(
		"%s?key=%s&channelId=%s&part=snippet,id&order=date&maxResults=%d",
		youtubeApiUrl,
		c.apiKey,
		channelId,
		maxResults,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var response YoutubeResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	var videos []VideoInfo
	for _, item := range response.Items {
		videos = append(
			videos, VideoInfo{
				Id:          item.Id.VideoId,
				Title:       item.Snippet.Title,
				PublishedAt: item.Snippet.PublishedAt,
			},
		)
	}

	return videos, nil
}

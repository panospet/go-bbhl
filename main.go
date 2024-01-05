package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kkdai/youtube/v2"

	"go-bbhl/util"
)

const (
	apiKey        = "AIzaSyDxFBI9kEBKLsWsX30ykHgoq7nByqEDCN4"
	channelID     = "UCGr3nR_XH9r6E5b09ZJAT9w"
	youtubeApiUrl = "https://www.googleapis.com/youtube/v3/search"
)

func main() {
	var dry bool
	flag.BoolVar(&dry, "dry", false, "dry run")
	flag.Parse()

	url := fmt.Sprintf(
		"%s?key=%s&channelId=%s&part=snippet,id&order=date&maxResults=50",
		youtubeApiUrl,
		apiKey,
		channelID,
	)

	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	var response YoutubeResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatalf("Error decoding response: %v", err)
	}

	var videos []VideoInfo
	for _, item := range response.Items {
		videos = append(
			videos, VideoInfo{
				Id:    item.Id.VideoId,
				Title: item.Snippet.Title,
			},
		)
	}

	// debug
	for _, v := range videos {
		log.Printf("Video: %v", v)
	}

	videos, err = filter(videos)
	if err != nil {
		log.Fatalf("Error filtering videos: %v", err)
	}

	for i, v := range videos {
		videoPath := fmt.Sprintf("./videos/%d.mp4", i)

		log.Printf("Downloading video: %s, id: %v, path: %s", v.Title, v.Id, videoPath)

		// create .videos directory if not exists
		if _, err := os.Stat("./videos"); os.IsNotExist(err) {
			if err := os.Mkdir("./videos", 0755); err != nil {
				log.Fatalf("Error creating directory: %v", err)
			}
		}

		if !dry {
			if err := downloadVideo(v.Id, videoPath); err != nil {
				log.Fatalf("Error downloading video: %v", err)
			}
		}
	}
}

func filter(videos []VideoInfo) ([]VideoInfo, error) {
	byRound := make(map[int][]VideoInfo)
	latestRound := 0
	for _, v := range videos {
		if !strings.Contains(strings.ToLower(v.Title), "highlights") {
			continue
		}

		round, err := util.ExtractElRound(v.Title)
		if err != nil {
			return nil, err
		}

		byRound[round] = append(byRound[round], v)
		if round > latestRound {
			latestRound = round
		}
	}

	log.Printf("Latest round: %v. Total videos: %d", latestRound, len(byRound[latestRound]))

	return byRound[latestRound], nil
}

type VideoInfo struct {
	Id    string
	Title string
}

func downloadVideo(
	videoId string,
	path string,
) error {
	client := youtube.Client{}

	url := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoId)

	videoInfo, err := client.GetVideo(url)
	if err != nil {
		return err
	}

	audioFormat := videoInfo.Formats.WithAudioChannels()
	videoStream, _, err := client.GetStream(videoInfo, &audioFormat[0])
	if err != nil {
		return err
	}

	videoFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer videoFile.Close()

	_, err = io.Copy(videoFile, videoStream)
	if err != nil {
		return err
	}

	log.Printf("Successfully downloaded: %v", url)
	return nil
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

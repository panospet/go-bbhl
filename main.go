package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"go-bbhl/downloader"
	"go-bbhl/util"
)

const (
	apiKey        = "AIzaSyDxFBI9kEBKLsWsX30ykHgoq7nByqEDCN4"
	channelID     = "UCGr3nR_XH9r6E5b09ZJAT9w"
	youtubeApiUrl = "https://www.googleapis.com/youtube/v3/search"
)

func main() {
	ctx := context.Background()

	var dry bool
	flag.BoolVar(&dry, "dry", false, "dry run")
	flag.Parse()

	url := fmt.Sprintf(
		"%s?key=%s&channelId=%s&part=snippet,id&order=date&maxResults=50",
		youtubeApiUrl,
		apiKey,
		channelID,
	)

	dl := downloader.NewYoutubeDownloader()

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

	videos, err = filter(videos)
	if err != nil {
		log.Fatalf("Error filtering videos: %v", err)
	}

	paths := make([]string, len(videos))
	for i, v := range videos {
		videoPath := fmt.Sprintf("./videos/%d.mp4", i)
		paths[i] = fmt.Sprintf("file '%s'", videoPath)

		log.Printf("Downloading video: %s, id: %v, path: %s", v.Title, v.Id, videoPath)

		// create .videos directory if not exists
		if _, err := os.Stat("./videos"); os.IsNotExist(err) {
			if err := os.Mkdir("./videos", 0755); err != nil {
				log.Fatalf("Error creating directory: %v", err)
			}
		}

		if !dry {
			if err := dl.DownloadVideo(v.Id, videoPath); err != nil {
				log.Fatalf("Error downloading video: %v", err)
			}
		}
	}

	if err := writeLinesToFile(paths, "./videos.txt"); err != nil {
		log.Fatalf("Error writing paths to file: %v", err)
	}

	// use ffmpeg terminal command to concat files
	// ffmpeg -f concat -safe 0 -i videos.txt -c copy output.mp4
	cmd := exec.Command(
		"ffmpeg", "-y", "-f", "concat", "-safe", "0", "-i", "videos.txt", "-c", "copy", "output.mp4",
	)

	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error running ffmpeg command: %v", err)
	}
	log.Printf("Successfully concatenated videos to output.mp4")

	// remove videos directory
	if err := os.RemoveAll("./videos"); err != nil {
		log.Fatalf("Error removing videos directory: %v", err)
	}
}

func writeLinesToFile(
	lines []string,
	path string,
) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return nil
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

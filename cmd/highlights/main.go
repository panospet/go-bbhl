package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/caarlos0/env/v9"

	"go-bbhl/downloader"
	"go-bbhl/filters"
	"go-bbhl/uploaders/telegram"
	"go-bbhl/youtube"
)

const (
	motionStationChannelId = "UCLd4dSmXdrJykO_hgOzbfPw"
	euroleagueChannelId    = "UCGr3nR_XH9r6E5b09ZJAT9w"
)

type config struct {
	TelegramToken string `env:"TELE_TOKEN,required"`
	ChatId        int64  `env:"CHAT_ID,required"`
	YoutubeApiKey string `env:"YOUTUBE_API_KEY,required"`
}

func main() {
	ctx := context.Background()

	var dry, nba, euroleague bool
	flag.BoolVar(&dry, "dry", false, "dry run")
	flag.BoolVar(&nba, "nba", false, "run for nba")
	flag.BoolVar(&euroleague, "euroleague", false, "run for euroleague")
	flag.Parse()

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	// should be either nba or euroleague
	if nba && euroleague || !nba && !euroleague {
		log.Fatalf("Only one of nba or euroleague flags should be set")
	}

	var channelId, caption string
	var filter func([]youtube.VideoInfo) ([]youtube.VideoInfo, error)
	if nba {
		channelId = motionStationChannelId
		filter = filters.NbaLatest
		caption = "latest NBA Highlights"
	}
	if euroleague {
		channelId = euroleagueChannelId
		filter = filters.EuroleagueLatestRound
		caption = "latest Euroleague Highlights"
	}

	start := time.Now()

	dl := downloader.NewYoutubeDownloader()

	ytClient := youtube.NewClient(cfg.YoutubeApiKey)

	upl, err := telegram.NewTelegramUploader(
		cfg.TelegramToken,
		cfg.ChatId,
		// todo change to env var
		"http://65.109.174.137:8081/bot%s/%s",
	)
	if err != nil {
		log.Fatalf("Error creating uploader: %v", err)
	}

	videos, err := ytClient.GetChannelVideos(
		ctx,
		channelId,
		50,
	)
	if err != nil {
		log.Fatalf("Error getting channel videos: %v", err)
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

	if dry {
		return
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

	log.Printf("uploading video to telegram channel...")
	if err := upl.UploadVideo("output.mp4", caption); err != nil {
		log.Fatalf("Error uploading video: %v", err)
	}

	// remove videos.txt and output.mp4
	if err := os.Remove("./videos.txt"); err != nil {
		log.Fatalf("Error removing videos.txt: %v", err)
	}
	if err := os.Remove("./output.mp4"); err != nil {
		log.Fatalf("Error removing output.mp4: %v", err)
	}

	log.Printf("finished after %v", time.Since(start))
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

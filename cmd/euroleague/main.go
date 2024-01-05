package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"go-bbhl/downloader"
	"go-bbhl/filters"
	"go-bbhl/youtube"
)

const (
	euroleagueChannelId = "UCGr3nR_XH9r6E5b09ZJAT9w"
)

func main() {
	ctx := context.Background()

	var dry bool
	flag.BoolVar(&dry, "dry", false, "dry run")
	flag.Parse()

	dl := downloader.NewYoutubeDownloader()

	ytClient := youtube.NewClient(os.Getenv("YOUTUBE_API_KEY"))

	videos, err := ytClient.GetChannelVideos(
		ctx,
		euroleagueChannelId,
		50,
	)
	if err != nil {
		log.Fatalf("Error getting channel videos: %v", err)
	}

	videos, err = filters.EuroleagueLatestRound(videos)
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

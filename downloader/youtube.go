package downloader

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/kkdai/youtube/v2"
)

type YoutubeDownloader struct {
	client *youtube.Client
}

func NewYoutubeDownloader() *YoutubeDownloader {
	return &YoutubeDownloader{
		client: &youtube.Client{},
	}
}

func (d *YoutubeDownloader) DownloadVideo(
	videoId string,
	path string,
) error {
	url := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoId)

	videoInfo, err := d.client.GetVideo(url)
	if err != nil {
		return err
	}

	audioFormat := videoInfo.Formats.WithAudioChannels()
	videoStream, _, err := d.client.GetStream(videoInfo, &audioFormat[0])
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

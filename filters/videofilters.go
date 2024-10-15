package filters

import (
	"log"
	"strings"
	"time"

	"go-bbhl/youtube"
)

func Euroleague(
	videos []youtube.VideoInfo,
) ([]youtube.VideoInfo, error) {
	var res []youtube.VideoInfo
	for _, v := range videos {
		if !strings.Contains(strings.ToLower(v.Title), "highlights") {
			continue
		}
		res = append(res, v)
	}

	log.Printf("Euroleague total videos: %d", len(res))

	return res, nil
}

func NbaLatest(
	videos []youtube.VideoInfo,
) ([]youtube.VideoInfo, error) {
	var res []youtube.VideoInfo
	for i, v := range videos {
		// if publishedAt is earlier than 10 hours before the previous video, break
		// that's how we assumed that game day is over
		if i > 0 {
			prevVideo := videos[i-1]
			if v.PublishedAt.Before(prevVideo.PublishedAt.Add(-10 * time.Hour)) {
				break
			}
		}

		res = append(res, v)
	}

	return res, nil
}

package filters

import (
	"log"
	"strings"
	"time"

	"go-bbhl/util"
	"go-bbhl/youtube"
)

func EuroleagueLatestRound(
	videos []youtube.VideoInfo,
) ([]youtube.VideoInfo, error) {
	byRound := make(map[int][]youtube.VideoInfo)
	latestRound := 0
	for _, v := range videos {
		if !strings.Contains(strings.ToLower(v.Title), "highlights") {
			continue
		}

		log.Printf("Processing video: %v", v.Title)

		// get last euroleague round based on video title
		round, err := util.ExtractElRound(v.Title)
		if err != nil {
			continue
		}

		byRound[round] = append(byRound[round], v)
		if round > latestRound {
			latestRound = round
		}
	}

	log.Printf("Latest round: %v. Total videos: %d", latestRound, len(byRound[latestRound]))

	return byRound[latestRound], nil
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

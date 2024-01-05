package filters

import (
	"log"
	"strings"

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

		// get last euroleague round based on video title
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

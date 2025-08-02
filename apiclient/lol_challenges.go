package apiclient

import (
	"fmt"

	"github.com/junioryono/Riot-API-Golang/apiclient/ratelimiter"
	"github.com/junioryono/Riot-API-Golang/constants/queue"
	"github.com/junioryono/Riot-API-Golang/constants/region"
)

type ChallengesConfig struct {
	ChallengeID string `json:"challengeId"`
	Levels      []struct {
		LevelID string `json:"levelId"`
		Order   int    `json:"order"`
	} `json:"levels"`
	QueueID queue.ID `json:"queueId"`
}

type ChallengesLeaderboards struct {
	Entries []struct {
		Entries []struct {
			Percentile float64 `json:"percentile"`
			Rank       int     `json:"rank"`
			Total      int     `json:"total"`
		} `json:"entries"`
		LevelID string `json:"levelId"`
	} `json:"entries"`
	QueueID queue.ID `json:"queueId"`
}

type ChallengesPercentiles struct {
	Percentiles []struct {
		Percentile float64 `json:"percentile"`
		Total      int     `json:"total"`
	} `json:"percentiles"`
	QueueID queue.ID `json:"queueId"`
}

type ChallengesPlayerData struct {
	Entries []struct {
		Entries []struct {
			Percentile float64 `json:"percentile"`
			Rank       int     `json:"rank"`
			Total      int     `json:"total"`
		} `json:"entries"`
		LevelID string `json:"levelId"`
	} `json:"entries"`
	QueueID queue.ID `json:"queueId"`
}

func (c *uniqueClient) GetChallengesConfig(r region.Region) (*ChallengesConfig, error) {
	var res ChallengesConfig
	err := c.dispatchAndUnmarshal(r, "/lol/challenges/v1/config", "", nil, ratelimiter.GetChallengesConfig, &res)
	return &res, err
}

func (c *uniqueClient) GetChallengesPercentiles(r region.Region) (*ChallengesPercentiles, error) {
	var res ChallengesPercentiles
	err := c.dispatchAndUnmarshal(r, "/lol/challenges/v1/percentiles", "", nil, ratelimiter.GetChallengesPercentiles, &res)
	return &res, err
}

func (c *uniqueClient) GetChallengesConfigByID(r region.Region, challengeID string) (*ChallengesConfig, error) {
	var res ChallengesConfig
	err := c.dispatchAndUnmarshal(r, "/lol/challenges/v1/config", fmt.Sprintf("/%s", challengeID), nil, ratelimiter.GetChallengesConfigByID, &res)
	return &res, err
}

func (c *uniqueClient) GetChallengesLeaderboardsByLevel(r region.Region, challengeID, level string) (*ChallengesLeaderboards, error) {
	var res ChallengesLeaderboards
	err := c.dispatchAndUnmarshal(r, "/lol/challenges/v1/leaderboards", fmt.Sprintf("/%s/%s", challengeID, level), nil, ratelimiter.GetChallengesLeaderboardsByLevel, &res)
	return &res, err
}

func (c *uniqueClient) GetChallengesPercentilesByID(r region.Region, challengeID string) (*ChallengesPercentiles, error) {
	var res ChallengesPercentiles
	err := c.dispatchAndUnmarshal(r, "/lol/challenges/v1/percentiles", fmt.Sprintf("/%s", challengeID), nil, ratelimiter.GetChallengesPercentilesByID, &res)
	return &res, err
}

func (c *uniqueClient) GetChallengesPlayerDataByPuuid(r region.Region, puuid string) (*ChallengesPlayerData, error) {
	var res ChallengesPlayerData
	err := c.dispatchAndUnmarshal(r, "/lol/challenges/v1/player-data", fmt.Sprintf("/%s", puuid), nil, ratelimiter.GetChallengesPlayerDataByPuuid, &res)
	return &res, err
}

package apiclient

import (
	"fmt"

	"github.com/junioryono/Riot-API-Golang/apiclient/ratelimiter"
	"github.com/junioryono/Riot-API-Golang/constants/region"
)

type ChampionMastery struct {
	Puuid                        string `json:"puuid"`
	ChampionID                   int32  `json:"championId"`
	ChampionLevel                int64  `json:"championLevel"`
	ChampionPoints               int64  `json:"championPoints"`
	LastPlayTime                 int64  `json:"lastPlayTime"`
	ChampionPointsSinceLastLevel int64  `json:"championPointsSinceLastLevel"`
	ChampionPointsUntilNextLevel int64  `json:"championPointsUntilNextLevel"`
	ChestGranted                 bool   `json:"chestGranted"`
	TokensEarned                 int64  `json:"tokensEarned"`
	SummonerID                   string `json:"summonerId"`
}

func (c *uniqueClient) GetChampionMasteriesBySummonerID(r region.Region, summonerID string) ([]ChampionMastery, error) {
	var res []ChampionMastery
	err := c.dispatchAndUnmarshal(r, "/lol/champion-mastery/v4/champion-masteries/by-summoner", fmt.Sprintf("/%s", summonerID), nil, ratelimiter.GetChampionMasteriesBySummonerID, &res)
	return res, err
}

func (c *uniqueClient) GetChampionMasteryBySummonerIDAndChampionID(r region.Region, summonerID string, championID int) (*ChampionMastery, error) {
	var res ChampionMastery
	err := c.dispatchAndUnmarshal(r, "/lol/champion-mastery/v4/champion-masteries/by-summoner", fmt.Sprintf("/%s/by-champion/%d", summonerID, championID), nil, ratelimiter.GetChampionMasteryBySummonerIDAndChampionID, &res)
	return &res, err
}

func (c *uniqueClient) GetChampionMasteriesTopBySummonerID(r region.Region, summonerID string) ([]ChampionMastery, error) {
	var res []ChampionMastery
	err := c.dispatchAndUnmarshal(r, "/lol/champion-mastery/v4/champion-masteries/by-summoner", fmt.Sprintf("/%s/top", summonerID), nil, ratelimiter.GetChampionMasteriesTopBySummonerID, &res)
	return res, err
}

func (c *uniqueClient) GetChampionMasteryScoreTotalBySummonerID(r region.Region, summonerID string) (int, error) {
	var res int
	err := c.dispatchAndUnmarshal(r, "/lol/champion-mastery/v4/scores/by-summoner", fmt.Sprintf("/%s", summonerID), nil, ratelimiter.GetChampionMasteryScoreTotalBySummonerID, &res)
	return res, err
}

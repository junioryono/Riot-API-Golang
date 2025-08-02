package apiclient

import (
	"fmt"

	"github.com/junioryono/Riot-API-Golang/apiclient/ratelimiter"
	"github.com/junioryono/Riot-API-Golang/constants/league/rank"
	"github.com/junioryono/Riot-API-Golang/constants/league/tier"
	"github.com/junioryono/Riot-API-Golang/constants/queue_ranked"
	"github.com/junioryono/Riot-API-Golang/constants/region"
)

type LeagueList struct {
	LeagueID string              `json:"leagueId"`
	Tier     tier.String         `json:"tier"`
	Entries  []LeagueItem        `json:"entries"`
	Queue    queue_ranked.String `json:"queue"`
	Name     string              `json:"name"`
}

type LeagueItem struct {
	SummonerID   string      `json:"summonerId"`
	Puuid        string      `json:"puuid"`
	LeaguePoints int32       `json:"leaguePoints"`
	Rank         rank.String `json:"rank"`
	Wins         int32       `json:"wins"`
	Losses       int32       `json:"losses"`
	Veteran      bool        `json:"veteran"`
	Inactive     bool        `json:"inactive"`
	FreshBlood   bool        `json:"freshBlood"`
	HotStreak    bool        `json:"hotStreak"`
}

type LeagueEntry struct {
	FreshBlood   bool                `json:"freshBlood"`
	HotStreak    bool                `json:"hotStreak"`
	Inactive     bool                `json:"inactive"`
	LeagueID     string              `json:"leagueId"`
	LeaguePoints int32               `json:"leaguePoints"`
	Losses       int32               `json:"losses"`
	QueueType    queue_ranked.String `json:"queueType"`
	Rank         rank.String         `json:"rank"`
	SummonerID   string              `json:"summonerId"`
	Puuid        string              `json:"puuid"`
	Tier         tier.String         `json:"tier"`
	Veteran      bool                `json:"veteran"`
	Wins         int32               `json:"wins"`
}

func (c *uniqueClient) GetLeagueEntriesChallenger(r region.Region, q queue_ranked.String) (*LeagueList, error) {
	var res LeagueList
	err := c.dispatchAndUnmarshal(r, "/lol/league/v4/challengerleagues/by-queue", fmt.Sprintf("/%s", q), nil, ratelimiter.GetLeagueEntriesChallenger, &res)
	return &res, err
}

func (c *uniqueClient) GetLeagueEntriesGrandmaster(r region.Region, q queue_ranked.String) (*LeagueList, error) {
	var res LeagueList
	err := c.dispatchAndUnmarshal(r, "/lol/league/v4/grandmasterleagues/by-queue", fmt.Sprintf("/%s", q), nil, ratelimiter.GetLeagueEntriesGrandmaster, &res)
	return &res, err
}

func (c *uniqueClient) GetLeagueEntriesMaster(r region.Region, q queue_ranked.String) (*LeagueList, error) {
	var res LeagueList
	err := c.dispatchAndUnmarshal(r, "/lol/league/v4/masterleagues/by-queue", fmt.Sprintf("/%s", q), nil, ratelimiter.GetLeagueEntriesMaster, &res)
	return &res, err
}

func (c *uniqueClient) GetLeagueEntries(r region.Region, q queue_ranked.String, tier tier.String, rank rank.String, page int) ([]LeagueEntry, error) {
	var res []LeagueEntry
	err := c.dispatchAndUnmarshal(r, "/lol/league/v4/entries", fmt.Sprintf("/%s/%s/%s?page=%d", q, tier, rank, page), nil, ratelimiter.GetLeagueEntries, &res)
	return res, err
}

func (c *uniqueClient) GetLeagueEntriesByID(r region.Region, leagueID string) (*LeagueList, error) {
	var res LeagueList
	err := c.dispatchAndUnmarshal(r, "/lol/league/v4/leagues", fmt.Sprintf("/%s", leagueID), nil, ratelimiter.GetLeagueEntriesByID, &res)
	return &res, err
}

func (c *uniqueClient) GetLeagueEntriesBySummonerID(r region.Region, summonerID string) ([]LeagueEntry, error) {
	var res []LeagueEntry
	err := c.dispatchAndUnmarshal(r, "/lol/league/v4/entries/by-summoner", fmt.Sprintf("/%s", summonerID), nil, ratelimiter.GetLeagueEntriesBySummonerID, &res)
	return res, err
}

func (c *uniqueClient) GetLeagueEntriesByPuuid(r region.Region, puuid string) ([]LeagueEntry, error) {
	var res []LeagueEntry
	err := c.dispatchAndUnmarshal(r, "/lol/league/v4/entries/by-puuid", fmt.Sprintf("/%s", puuid), nil, ratelimiter.GetLeagueEntriesByPuuid, &res)
	return res, err
}

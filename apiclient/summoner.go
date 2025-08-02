package apiclient

import (
	"encoding/json"
	"fmt"

	"github.com/junioryono/Riot-API-Golang/apiclient/ratelimiter"
	"github.com/junioryono/Riot-API-Golang/constants/region"
)

type Summoner struct {
	AccountID     string `json:"accountId"`
	ID            string `json:"id"`
	ProfileIconID int32  `json:"profileIconId"`
	Puuid         string `json:"puuid"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int32  `json:"summonerLevel"`
}

func (s Summoner) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Summoner) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

func (c *uniqueClient) GetSummonerByRsoPuuid(r region.Region, rsoPuuid string) (*Summoner, error) {
	var res Summoner
	err := c.dispatchAndUnmarshal(r, "/lol/summoner/v4/summoners/by-account", fmt.Sprintf("/%s", rsoPuuid), nil, ratelimiter.GetSummonerByRsoPuuid, &res)
	return &res, err
}

func (c *uniqueClient) GetSummonerByAccountID(r region.Region, accountID string) (*Summoner, error) {
	var res Summoner
	err := c.dispatchAndUnmarshal(r, "/lol/summoner/v4/summoners/by-account", fmt.Sprintf("/%s", accountID), nil, ratelimiter.GetSummonerByAccountID, &res)
	return &res, err
}

func (c *uniqueClient) GetSummonerByPuuid(r region.Region, puuid string) (*Summoner, error) {
	var res Summoner
	err := c.dispatchAndUnmarshal(r, "/lol/summoner/v4/summoners/by-puuid", fmt.Sprintf("/%s", puuid), nil, ratelimiter.GetSummonerByPuuid, &res)
	return &res, err
}

func (c *uniqueClient) GetSummonerBySummonerID(r region.Region, summonerID string) (*Summoner, error) {
	var res Summoner
	err := c.dispatchAndUnmarshal(r, "/lol/summoner/v4/summoners", fmt.Sprintf("/%s", summonerID), nil, ratelimiter.GetSummonerBySummonerID, &res)
	return &res, err
}

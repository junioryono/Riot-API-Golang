package apiclient

import (
	"fmt"

	"github.com/junioryono/Riot-API-Golang/apiclient/ratelimiter"
	"github.com/junioryono/Riot-API-Golang/constants/continent"
)

type Account struct {
	Puuid    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

func (c *uniqueClient) GetAccountByPuuid(continent continent.Continent, puuid string) (*Account, error) {
	var account Account
	err := c.dispatchAndUnmarshal(continent, "/riot/account/v1/accounts/by-puuid", fmt.Sprintf("/%s", puuid), nil, ratelimiter.GetAccountByPuuid, &account)
	return &account, err
}

func (c *uniqueClient) GetAccountByRiotID(continent continent.Continent, gameName, tagLine string) (*Account, error) {
	var account Account
	err := c.dispatchAndUnmarshal(continent, "/riot/account/v1/accounts/by-riot-id", fmt.Sprintf("/%s/%s", gameName, tagLine), nil, ratelimiter.GetAccountByRiotID, &account)
	return &account, err
}

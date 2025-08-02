package apiclient

import (
	"github.com/junioryono/Riot-API-Golang/apiclient/ratelimiter"
	"github.com/junioryono/Riot-API-Golang/constants/region"
)

type ChampionRotations struct {
	FreeChampionIDs              []int16 `json:"freeChampionIds"`
	FreeChampionIDsForNewPlayers []int16 `json:"freeChampionIdsForNewPlayers"`
	MaxNewPlayerLevel            int16   `json:"maxNewPlayerLevel"`
}

func (c *uniqueClient) GetChampionRotations(r region.Region) (*ChampionRotations, error) {
	var championRotation ChampionRotations
	err := c.dispatchAndUnmarshal(r, "/lol/platform/v3/champion-rotations", "", nil, ratelimiter.GetChampionRotations, &championRotation)
	return &championRotation, err
}

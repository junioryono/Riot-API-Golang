package apiclient

import (
	"fmt"

	"github.com/junioryono/Riot-API-Golang/apiclient/ratelimiter"
	"github.com/junioryono/Riot-API-Golang/constants/queue"
	"github.com/junioryono/Riot-API-Golang/constants/region"
	"github.com/junioryono/Riot-API-Golang/constants/summoner_spell"
)

type ActiveGame struct {
	GameID            int64                   `json:"gameId"`            // The ID of the game
	MapID             int16                   `json:"mapId"`             // The ID of the map
	GameMode          string                  `json:"gameMode"`          // The game mode
	GameType          string                  `json:"gameType"`          // The game type
	GameQueueConfigID queue.ID                `json:"gameQueueConfigId"` // The queue type (queue types are documented on the Game Constants page)
	Participants      []ActiveGameParticipant `json:"participants"`      // The participant information
	Observers         Observers               `json:"observers"`         // The observer information
	PlatformID        region.Region           `json:"platformId"`        // The ID of the platform on which the game is being played
	BannedChampions   []BannedChampion        `json:"bannedChampions"`   // Banned champion information
	GameStartTime     int64                   `json:"gameStartTime"`     // The game start time represented in epoch milliseconds
	GameLength        int64                   `json:"gameLength"`        // The amount of time in seconds that has passed since the game started
}

type ActiveGameParticipant struct {
	Puuid         *string           `json:"puuid"`         // The Puuid of the player
	RiotID        *string           `json:"riotId"`        // The Riot ID of the player
	SummonerID    *string           `json:"summonerId"`    // The encrypted summoner ID of this participant
	ProfileIconID *int32            `json:"profileIconId"` // The ID of the profile icon used by this participant
	ChampionID    int32             `json:"championId"`    // The ID of the champion played by this participant
	TeamID        int16             `json:"teamId"`        // The team ID of this participant, indicating the participant's team
	Bot           bool              `json:"bot"`           // Flag indicating whether or not this participant is a bot
	Spell1ID      summoner_spell.ID `json:"spell1Id"`      // The ID of the first summoner spell used by this participant
	Spell2ID      summoner_spell.ID `json:"spell2Id"`      // The ID of the second summoner spell used by this participant
	Perks         Perks             `json:"perks"`
}

type Perks struct {
	PerkIDs      []int16 `json:"perkIds"`
	PerkStyle    int16   `json:"perkStyle"`
	PerkSubStyle int16   `json:"perkSubStyle"`
}

type Observers struct {
	EncryptionKey string `json:"encryptionKey"` // Key used to decrypt the spectator grid game data for playback
}

type BannedChampion struct {
	ChampionID int32 `json:"championId"` // The ID of the banned champion
	TeamID     int16 `json:"teamId"`     // The ID of the team that banned the champion
	PickTurn   int16 `json:"pickTurn"`   // The turn during which the champion was banned
}

func (c *uniqueClient) GetSpectatorActiveGameByPuuid(r region.Region, puuid string) (*ActiveGame, error) {
	var res ActiveGame
	err := c.dispatchAndUnmarshal(r, "/lol/spectator/v5/active-games/by-summoner", fmt.Sprintf("/%s", puuid), nil, ratelimiter.GetSpectatorActiveGameByPuuid, &res)
	return &res, err
}

type FeaturedGames struct {
	GameList []FeaturedGame `json:"gameList"` // The list of featured games
}

type FeaturedGame struct {
	GameID            int64                   `json:"gameId"`            // The ID of the game
	MapID             int16                   `json:"mapId"`             // The ID of the map
	GameMode          string                  `json:"gameMode"`          // The game mode
	GameType          string                  `json:"gameType"`          // The game type
	GameQueueConfigID queue.ID                `json:"gameQueueConfigId"` // The queue type (queue types are documented on the Game Constants page)
	Participants      []ActiveGameParticipant `json:"participants"`      // The participant information
	Observers         Observers               `json:"observers"`         // The observer information
	PlatformID        region.Region           `json:"platformId"`        // The ID of the platform on which the game is being played
	BannedChampions   []BannedChampion        `json:"bannedChampions"`   // Banned champion information
	GameStartTime     int64                   `json:"gameStartTime"`     // The game start time represented in epoch milliseconds
	GameLength        int64                   `json:"gameLength"`        // The amount of time in seconds that has passed since the game started
}

type FeaturedGameParticipant struct {
	Puuid         *string           `json:"puuid"`         // The Puuid of the player
	RiotID        *string           `json:"riotId"`        // The Riot ID of the player
	SummonerID    *string           `json:"summonerId"`    // The encrypted summoner ID of this participant
	ChampionID    int32             `json:"championId"`    // The ID of the champion played by this participant
	TeamID        int16             `json:"teamId"`        // The team ID of this participant, indicating the participant's team
	ProfileIconID int32             `json:"profileIconId"` // The ID of the profile icon used by this participant
	Bot           bool              `json:"bot"`           // Flag indicating whether or not this participant is a bot
	Spell1ID      summoner_spell.ID `json:"spell1Id"`      // The ID of the first summoner spell used by this participant
	Spell2ID      summoner_spell.ID `json:"spell2Id"`      // The ID of the second summoner spell used by this participant
	Perks         Perks             `json:"perks"`
}

func (c *uniqueClient) GetSpectatorFeaturedGames(r region.Region) (*FeaturedGames, error) {
	var res FeaturedGames
	err := c.dispatchAndUnmarshal(r, "/lol/spectator/v5/featured-games", "", nil, ratelimiter.GetSpectatorFeaturedGames, &res)
	return &res, err
}

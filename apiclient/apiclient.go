package apiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/junioryono/Riot-API-Golang/apiclient/ratelimiter"
	"github.com/junioryono/Riot-API-Golang/constants/continent"
	"github.com/junioryono/Riot-API-Golang/constants/league/rank"
	"github.com/junioryono/Riot-API-Golang/constants/league/tier"
	"github.com/junioryono/Riot-API-Golang/constants/queue_ranked"
	"github.com/junioryono/Riot-API-Golang/constants/region"
	"github.com/barkimedes/go-deepcopy"
)

type Client interface {
	// WithContext and WithPriority are used to set the context and priority of the request.
	WithContext(ctx context.Context) Client
	WithPriority(priority int) Client
	WithCache(duration time.Duration) Client

	// Helper methods to set the API key, usage conservation, and max retries.
	SetUsageConservation(conserveUsage ratelimiter.ConserveUsage)
	SetAPIKey(apiKey string)
	SetMaxRetries(maxRetries int)
	SetCacheCleanupDuration(duration time.Duration)

	/* Account API */

	GetAccountByPuuid(continent continent.Continent, puuid string) (*Account, error)
	GetAccountByRiotID(continent continent.Continent, gameName, tagLine string) (*Account, error)

	/* Champion Mastery API */

	GetChampionMasteriesBySummonerID(region region.Region, summonerID string) ([]ChampionMastery, error)
	GetChampionMasteryBySummonerIDAndChampionID(region region.Region, summonerID string, championID int) (*ChampionMastery, error)
	GetChampionMasteriesTopBySummonerID(region region.Region, summonerID string) ([]ChampionMastery, error)
	GetChampionMasteryScoreTotalBySummonerID(region region.Region, summonerID string) (int, error)

	/* Champion API */

	GetChampionRotations(region region.Region) (*ChampionRotations, error)

	/* Clash API */

	GetClashPlayersByPuuid(region region.Region, puuid string) (*ClashPlayers, error)
	GetClashPlayersBySummonerID(region region.Region, summonerID string) (*ClashPlayers, error)
	GetClashTeamByID(region region.Region, teamID string) (*ClashTeam, error)
	GetClashTournaments(region region.Region) (*ClashTournaments, error)
	GetClashTournamentByTeamID(region region.Region, teamID string) (*ClashTournament, error)
	GetClashTournamentByID(region region.Region, tournamentID string) (*ClashTournament, error)

	/* League Exp API */

	GetLeagueExpEntries(region region.Region, q queue_ranked.String, tier tier.String, rank rank.String, page int) ([]LeagueEntry, error)

	/* League API */

	GetLeagueEntriesChallenger(region region.Region, q queue_ranked.String) (*LeagueList, error)
	GetLeagueEntriesGrandmaster(region region.Region, q queue_ranked.String) (*LeagueList, error)
	GetLeagueEntriesMaster(region region.Region, q queue_ranked.String) (*LeagueList, error)
	GetLeagueEntries(region region.Region, q queue_ranked.String, tier tier.String, rank rank.String, page int) ([]LeagueEntry, error)
	GetLeagueEntriesByID(region region.Region, leagueID string) (*LeagueList, error)
	GetLeagueEntriesBySummonerID(region region.Region, summonerID string) ([]LeagueEntry, error)
	GetLeagueEntriesByPuuid(region region.Region, puuid string) ([]LeagueEntry, error)

	/* LOL Challenges API */

	GetChallengesConfig(region region.Region) (*ChallengesConfig, error)
	GetChallengesPercentiles(region region.Region) (*ChallengesPercentiles, error)
	GetChallengesConfigByID(region region.Region, challengeID string) (*ChallengesConfig, error)
	GetChallengesLeaderboardsByLevel(region region.Region, challengeID, level string) (*ChallengesLeaderboards, error)
	GetChallengesPercentilesByID(region region.Region, challengeID string) (*ChallengesPercentiles, error)
	GetChallengesPlayerDataByPuuid(region region.Region, puuid string) (*ChallengesPlayerData, error)

	/* LOL Status API */

	GetStatusPlatformData(region region.Region) (*StatusPlatformData, error)

	/* Match API */

	GetMatchlist(continent continent.Continent, puuid string, opts *GetMatchlistOptions) (*Matchlist, error)
	GetMatch(continent continent.Continent, matchID string) (*Match, error)
	GetMatchTimeline(continent continent.Continent, matchID string) (*MatchTimeline, error)

	/* Spectator API */

	GetSpectatorActiveGameByPuuid(region region.Region, summonerID string) (*ActiveGame, error)
	GetSpectatorFeaturedGames(region region.Region) (*FeaturedGames, error)

	/* Summoner API */

	GetSummonerByRsoPuuid(region region.Region, rsoPuuid string) (*Summoner, error)
	GetSummonerByAccountID(region region.Region, accountID string) (*Summoner, error)
	GetSummonerByPuuid(region region.Region, puuid string) (*Summoner, error)
	GetSummonerBySummonerID(region region.Region, summonerID string) (*Summoner, error)
}

type cacheEntry struct {
	data   interface{}
	expiry time.Time
}

type sharedClient struct {
	ratelimiter          *ratelimiter.RateLimiter
	cache                map[string]*cacheEntry
	cacheMutex           sync.Mutex
	cacheCleanupDuration time.Duration
}

type uniqueClient struct {
	*sharedClient
	ctx           context.Context
	priority      int
	cacheDuration time.Duration
}

func New(apiKey string) Client {
	requests := make(chan *ratelimiter.APIRequest)
	ratelimiter := ratelimiter.NewRateLimiter(requests, apiKey)
	go ratelimiter.Start()

	u := &sharedClient{
		ratelimiter:          ratelimiter,
		cache:                make(map[string]*cacheEntry),
		cacheCleanupDuration: 5 * time.Minute,
	}

	c := &uniqueClient{
		sharedClient: u,
		ctx:          context.Background(),
	}

	go c.cleanupCache()

	return c
}

func (c *uniqueClient) WithContext(ctx context.Context) Client {
	return &uniqueClient{
		sharedClient:  c.sharedClient,
		ctx:           ctx,
		priority:      c.priority,
		cacheDuration: c.cacheDuration,
	}
}

func (c *uniqueClient) WithPriority(priority int) Client {
	return &uniqueClient{
		sharedClient:  c.sharedClient,
		ctx:           c.ctx,
		priority:      priority,
		cacheDuration: c.cacheDuration,
	}
}

func (c *uniqueClient) WithCache(duration time.Duration) Client {
	return &uniqueClient{
		sharedClient:  c.sharedClient,
		ctx:           c.ctx,
		priority:      c.priority,
		cacheDuration: duration,
	}
}

func (c *uniqueClient) SetUsageConservation(conserveUsage ratelimiter.ConserveUsage) {
	c.ratelimiter.SetUsageConservation(conserveUsage)
}

func (c *uniqueClient) SetAPIKey(apiKey string) {
	c.ratelimiter.SetAPIKey(apiKey)
}

func (c *uniqueClient) SetMaxRetries(maxRetries int) {
	c.ratelimiter.SetMaxRetries(maxRetries)
}

func (c *uniqueClient) SetCacheCleanupDuration(duration time.Duration) {
	c.cacheCleanupDuration = duration
}

type HostProvider interface {
	Host() string
	String() string
}

func (c *uniqueClient) dispatchAndUnmarshal(regionOrContinent HostProvider, method string, relativePath string, parameters url.Values, methodID ratelimiter.MethodID, dest interface{}) error {
	var suffix, separator string

	if len(parameters) > 0 {
		suffix = fmt.Sprintf("?%s", parameters.Encode())
	}

	if !strings.HasPrefix(relativePath, "/") {
		separator = "/"
	}

	URL := regionOrContinent.Host() + method + separator + relativePath + suffix

	// Check if in cache
	if cachedData, ok := c.getFromCache(URL); ok {
		copy, err := deepcopy.Anything(cachedData)
		if err != nil {
			return err
		}

		destValue := reflect.ValueOf(dest)
		if destValue.Kind() != reflect.Ptr {
			return fmt.Errorf("destination must be a pointer")
		}
		destElem := destValue.Elem()

		copyValue := reflect.ValueOf(copy)
		if copyValue.Kind() == reflect.Ptr {
			copyValue = copyValue.Elem()
		}

		if !copyValue.Type().AssignableTo(destElem.Type()) {
			return fmt.Errorf("cached data type %v is not assignable to destination type %v", copyValue.Type(), destElem.Type())
		}

		destElem.Set(copyValue)

		return nil
	}

	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)
	newRequest := ratelimiter.APIRequest{
		Context:  c.ctx,
		Priority: c.priority,
		Region:   strings.ToUpper(regionOrContinent.String()),
		MethodID: methodID,
		URL:      URL,
		Response: responseChan,
		Error:    errorChan,
	}

	// Insert the request into the rate limiter
	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
	case c.ratelimiter.Requests <- &newRequest:
	}

	// Wait for the response
	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
	case err := <-errorChan:
		return err
	case response := <-responseChan:
		if response == nil {
			return fmt.Errorf("received nil response (%s)", URL)
		}

		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			if err, ok := StatusToError[response.StatusCode]; ok {
				return fmt.Errorf("status code %d: %s (%s)", response.StatusCode, err, URL)
			}

			return fmt.Errorf("status code %d: unknown error (%s)", response.StatusCode, URL)
		}

		decoder := json.NewDecoder(response.Body)
		if err := decoder.Decode(dest); err != nil {
			return fmt.Errorf("failed to decode response: %w (%s)", err, URL)
		}

		// Cache the destination
		if c.cacheDuration > 0 {
			c.addToCache(URL, dest)
		}

		return nil
	}
}

func (c *uniqueClient) getFromCache(URL string) (interface{}, bool) {
	c.cacheMutex.Lock()

	if entry, ok := c.cache[URL]; ok {
		if time.Now().Before(entry.expiry) {
			value := entry.data
			c.cacheMutex.Unlock()
			return value, true
		}

		delete(c.cache, URL)
	}

	c.cacheMutex.Unlock()
	return nil, false
}

func (c *uniqueClient) addToCache(URL string, data interface{}) {
	copy, err := deepcopy.Anything(data)
	if err != nil {
		return
	}

	c.cacheMutex.Lock()
	c.cache[URL] = &cacheEntry{
		data:   copy,
		expiry: time.Now().Add(c.cacheDuration),
	}
	c.cacheMutex.Unlock()
}

func (c *uniqueClient) cleanupCache() {
	ticker := time.NewTicker(c.cacheCleanupDuration)
	defer ticker.Stop()

	for {
		<-ticker.C
		c.removeExpiredEntries()
	}
}

func (c *uniqueClient) removeExpiredEntries() {
	c.cacheMutex.Lock()

	now := time.Now()
	for URL, entry := range c.cache {
		if now.After(entry.expiry) {
			delete(c.cache, URL)
		}
	}

	c.cacheMutex.Unlock()
}

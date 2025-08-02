[![Automated Version Tagging](https://github.com/junioryono/Riot-API-Golang/actions/workflows/tagging.yml/badge.svg)](https://github.com/junioryono/Riot-API-Golang/actions/workflows/tagging.yml)

# ⚡️ Riot-API-Golang

> A sleek, easy-to-use wrapper for making API calls to Riot and DataDragon with built-in rate limiting

The Riot-API-Golang provides a simplified interface to access Riot Games and DataDragon APIs in Golang. It smoothly handles Riot's rate limiting, making your development experience hassle-free.

## Features

- 🚀 Simple to use
- 🧠 Intelligent rate limiting
- 🎮 Access to Riot Games and DataDragon APIs

## Installation

```bash
go get github.com/junioryono/Riot-API-Golang
```

## Example Usage (Riot API)

continent continent.Continent, gameName, tagLine string

```go
func main() {
	apiKey := "RGAPI-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	client := apiclient.New(apiKey)

	riotAccount, err := client.GetAccountByRiotID(continent.AMERICAS, "Mighty Junior", "NA1")
	if err != nil {
		panic(err)
	}

	matchlist, err := client.GetMatchlist(region.NA1.Continent(), riotAccount.Puuid, nil)
	if err != nil {
		panic(err)
	}

	matchID := (*matchlist)[0]

	match, err := client.GetMatch(region.NA1.Continent(), matchID)
	if err != nil {
		panic(err)
	}

	fmt.Println(match.Info.GameCreation)
}
```

## Example Usage (Riot API - With Context)

```go
func main() {
	apiKey := "RGAPI-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	client := apiclient.New(apiKey)

	riotAccount, err := client.GetAccountByRiotID(continent.AMERICAS, "Mighty Junior", "NA1")
	if err != nil {
		panic(err)
	}

	matchlist, err := client.WithContext(context.TODO()).GetMatchlist(region.NA1.Continent(), riotAccount.Puuid, nil)
	if err != nil {
		panic(err)
	}

	matchID := (*matchlist)[0]

	match, err := client.WithContext(context.TODO()).GetMatch(region.NA1.Continent(), matchID)
	if err != nil {
		panic(err)
	}

	fmt.Println(match.Info.GameCreation)
}
```

## Example Usage (DDragon)

```go
func main() {
	patches, err := staticdata.GetPatches()
	if err != nil {
		panic(err)
	}

	currentPatch := patches.CurrentPatch()

	champions, err := staticdata.GetChampions(currentPatch, language.EnglishUnitedStates)
	if err != nil {
		panic(err)
	}

	champ, err := champions.Champion("MonkeyKing")
	if err != nil {
		panic(err)
	}

	fmt.Println(champ.Blurb)
}
```

## Request Throttling

Throttle the number of requests made to Riot's APIs.
This is beneficial if a frontend portion of the application exists.

```go
conservation := ratelimiter.ConserveUsage{
    RegionPercent: 30,
    MethodPercent: 30,
}

client.SetUsageConservation(conservation)
```

Ignore limits for specific methods.
The region's conservation percentage will still be followed.

```go
conservation := ratelimiter.ConserveUsage{
    RegionPercent: 30,
    MethodPercent: 30,
    IgnoreLimits: []ratelimiter.MethodID{
        ratelimiter.GetLeagueEntriesChallenger,
        ratelimiter.GetLeagueEntriesGrandmaster,
        ratelimiter.GetLeagueEntriesMaster,
        ratelimiter.GetLeagueEntries,
    },
}

client.SetUsageConservation(conservation)
```

## Request Error Handling

How many times Riot API requests will be retried when unsuccessful. By default, requests will be retried indefinitely (-1).

```go
client.SetMaxRetries(3)
```

## Contributing

Interested in contributing to Riot-API-Golang? Check out the [contributing guide](CONTRIBUTING.md) to see how you can make an impact.

## License

Riot-API-Golang is licensed under the MIT license. See the [LICENSE](LICENSE) file for more info.

## Support

If you encounter any issues or have questions, please file an issue on the [GitHub issues page](https://github.com/junioryono/Riot-API-Golang/issues).

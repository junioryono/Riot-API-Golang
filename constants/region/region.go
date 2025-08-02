package region

import (
	"strings"
	"unicode"

	"github.com/junioryono/Riot-API-Golang/constants/continent"
)

type Region string

const (
	// Brazil
	BR1 Region = "BR1"

	// Europe East
	EUN1 Region = "EUN1"

	// Europe West
	EUW1 Region = "EUW1"

	// Japan
	JP1 Region = "JP1"

	// Korea
	KR Region = "KR"

	// Latin America North
	LA1 Region = "LA1"

	// Latin America South
	LA2 Region = "LA2"

	// Middle East
	ME1 Region = "ME1"

	// North America
	NA1 Region = "NA1"

	// Oceania
	OC1 Region = "OC1"

	// Russia
	RU Region = "RU"

	// Southeast Asia
	SEA Region = "SG2"

	// Turkey
	TR1 Region = "TR1"

	// Taiwan
	TW2 Region = "TW2"

	// Vietnam
	VN2 Region = "VN2"
)

func (r Region) String() string {
	return string(r)
}

var stringToRegion = map[string]Region{
	"BR":  BR1,
	"EUN": EUN1,
	"EUW": EUW1,
	"JP":  JP1,
	"KR":  KR,
	"LA1": LA1,
	"LAN": LA1,
	"LA2": LA2,
	"LAS": LA2,
	"ME":  ME1,
	"NA":  NA1,
	"OC":  OC1,
	"RU":  RU,
	"SEA": SEA,
	"SG":  SEA,
	"TR":  TR1,
	"TW":  TW2,
	"VN":  VN2,
}

func FromString(rgn string) (Region, bool) {
	// Capitalize the string
	rgn = strings.ToUpper(rgn)

	// Special handling for Latin America regions
	if strings.HasPrefix(rgn, "LA") {
		if region, ok := stringToRegion[rgn]; ok {
			return region, true
		}
	} else {
		// Remove all numbers from the string for other regions
		rgn = strings.Map(func(r rune) rune {
			if unicode.IsNumber(r) {
				return -1
			}
			return r
		}, rgn)
	}

	region, ok := stringToRegion[rgn]
	return region, ok
}

var regionToHost = map[Region]string{
	BR1:  "https://br1.api.riotgames.com",
	EUN1: "https://eun1.api.riotgames.com",
	EUW1: "https://euw1.api.riotgames.com",
	JP1:  "https://jp1.api.riotgames.com",
	KR:   "https://kr.api.riotgames.com",
	LA1:  "https://la1.api.riotgames.com",
	LA2:  "https://la2.api.riotgames.com",
	ME1:  "https://me1.api.riotgames.com",
	NA1:  "https://na1.api.riotgames.com",
	OC1:  "https://oc1.api.riotgames.com",
	RU:   "https://ru.api.riotgames.com",
	SEA:  "https://sg2.api.riotgames.com",
	TR1:  "https://tr1.api.riotgames.com",
	TW2:  "https://tw2.api.riotgames.com",
	VN2:  "https://vn2.api.riotgames.com",
}

// Returns the full hostname corresponding to the region.
func (r Region) Host() string {
	return regionToHost[r]
}

var regionToContinentMatchV5 = map[Region]continent.Continent{
	BR1:  continent.AMERICAS,
	EUN1: continent.EUROPE,
	EUW1: continent.EUROPE,
	JP1:  continent.ASIA,
	KR:   continent.ASIA,
	LA1:  continent.AMERICAS,
	LA2:  continent.AMERICAS,
	ME1:  continent.EUROPE,
	NA1:  continent.AMERICAS,
	OC1:  continent.SEA,
	RU:   continent.EUROPE,
	SEA:  continent.SEA,
	TR1:  continent.EUROPE,
	TW2:  continent.SEA,
	VN2:  continent.SEA,
}

// Returns the continent that the region is in
func (r Region) ContinentMatchV5() continent.Continent {
	return regionToContinentMatchV5[r]
}

// Map the nearest region to the continent
var regionToContinentAccountV1 = map[Region]continent.Continent{
	BR1:  continent.AMERICAS,
	EUN1: continent.EUROPE,
	EUW1: continent.EUROPE,
	JP1:  continent.ASIA,
	KR:   continent.ASIA,
	LA1:  continent.AMERICAS,
	LA2:  continent.AMERICAS,
	ME1:  continent.EUROPE,
	NA1:  continent.AMERICAS,
	OC1:  continent.AMERICAS,
	RU:   continent.ASIA,
	SEA:  continent.ASIA,
	TR1:  continent.EUROPE,
	TW2:  continent.ASIA,
	VN2:  continent.ASIA,
}

// Returns the continent that the region is in
func (r Region) ContinentAccountV1() continent.Continent {
	return regionToContinentAccountV1[r]
}

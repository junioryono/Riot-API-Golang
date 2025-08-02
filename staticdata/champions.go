package staticdata

import (
	"fmt"

	"github.com/junioryono/Riot-API-Golang/constants/language"
	"github.com/junioryono/Riot-API-Golang/constants/patch"
)

type Champions []Champion

type Champion struct {
	ID      string        `json:"id"`
	Key     string        `json:"key"`
	Name    string        `json:"name"`
	Title   string        `json:"title"`
	Blurb   string        `json:"blurb"`
	Info    ChampionInfo  `json:"info"`
	Image   Image         `json:"image"`
	Tags    []string      `json:"tags"`
	Partype string        `json:"partype"`
	Stats   ChampionStats `json:"stats"`
}

type ChampionInfo struct {
	Difficulty int `json:"difficulty"`
	Attack     int `json:"attack"`
	Defense    int `json:"defense"`
	Magic      int `json:"magic"`
}

type ChampionStats struct {
	Armorperlevel        float64 `json:"armorperlevel"`
	Hpperlevel           float64 `json:"hpperlevel"`
	Attackdamage         float64 `json:"attackdamage"`
	Mpperlevel           float64 `json:"mpperlevel"`
	Attackspeedoffset    float64 `json:"attackspeedoffset"`
	Armor                float64 `json:"armor"`
	Hp                   float64 `json:"hp"`
	Hpregenperlevel      float64 `json:"hpregenperlevel"`
	Spellblock           float64 `json:"spellblock"`
	Attackrange          float64 `json:"attackrange"`
	Movespeed            float64 `json:"movespeed"`
	Attackdamageperlevel float64 `json:"attackdamageperlevel"`
	Mpregenperlevel      float64 `json:"mpregenperlevel"`
	Mp                   float64 `json:"mp"`
	Spellblockperlevel   float64 `json:"spellblockperlevel"`
	Crit                 float64 `json:"crit"`
	Mpregen              float64 `json:"mpregen"`
	Attackspeedperlevel  float64 `json:"attackspeedperlevel"`
	Hpregen              float64 `json:"hpregen"`
	Critperlevel         float64 `json:"critperlevel"`
}

func GetChampions(v patch.Patch, lang language.Language) (Champions, error) {
	type Response struct {
		Data map[string]Champion `json:"data"`
	}

	var res Response
	err := getJSON(fmt.Sprintf("http://ddragon.leagueoflegends.com/cdn/%s/data/%s/champion.json", v, lang), &res)

	var champions Champions
	for _, champion := range res.Data {
		champions = append(champions, champion)
	}

	return champions, err
}

func (c Champions) Champion(championID string) (Champion, error) {
	for _, champion := range c {
		if champion.ID == championID {
			return champion, nil
		}
	}

	return Champion{}, fmt.Errorf("champion %s not found", championID)
}

type ChampionDetailed struct {
	ID        string          `json:"id"`
	Key       string          `json:"key"`
	Name      string          `json:"name"`
	Title     string          `json:"title"`
	Image     Image           `json:"image"`
	Skins     []ChampionSkin  `json:"skins"`
	Lore      string          `json:"lore"`
	Blurb     string          `json:"blurb"`
	Allytips  []string        `json:"allytips"`
	Enemytips []string        `json:"enemytips"`
	Tags      []string        `json:"tags"`
	Partype   string          `json:"partype"`
	Info      ChampionInfo    `json:"info"`
	Stats     ChampionStats   `json:"stats"`
	Spells    []ChampionSpell `json:"spells"`
	Passive   ChampionPassive `json:"passive"`
}

type ChampionSkin struct {
	ID      string `json:"id"`
	Num     int    `json:"num"`
	Name    string `json:"name"`
	Chromas bool   `json:"chromas"`
}

type ChampionSpell struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Tooltip      string             `json:"tooltip"`
	Leveltip     ChampionLevelTip   `json:"leveltip"`
	Maxrank      int                `json:"maxrank"`
	Cooldown     []float64          `json:"cooldown"`
	CooldownBurn string             `json:"cooldownBurn"`
	Cost         []int              `json:"cost"`
	CostBurn     string             `json:"costBurn"`
	Effect       [][]float64        `json:"effect"`
	EffectBurn   []string           `json:"effectBurn"`
	Vars         []ChampionSpellVar `json:"vars"`
	CostType     string             `json:"costType"`
	MaxAmmo      string             `json:"maxammo"`
	Range        []int              `json:"range"`
	RangeBurn    string             `json:"rangeBurn"`
	Image        Image              `json:"image"`
	Resource     string             `json:"resource"`
}

type ChampionLevelTip struct {
	Label  []string `json:"label"`
	Effect []string `json:"effect"`
}

type ChampionSpellVar struct {
	Key   string    `json:"key"`
	Link  string    `json:"link"`
	Coeff []float64 `json:"coeff"`
}

type ChampionPassive struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       Image  `json:"image"`
}

func GetChampion(v patch.Patch, lang language.Language, championID string) (ChampionDetailed, error) {
	type Response struct {
		Data map[string]ChampionDetailed `json:"data"`
	}

	var res Response
	err := getJSON(fmt.Sprintf("http://ddragon.leagueoflegends.com/cdn/%s/data/%s/champion/%s.json", v, lang, championID), &res)
	return res.Data[championID], err
}

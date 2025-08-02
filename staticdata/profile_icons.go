package staticdata

import (
	"fmt"
	"strconv"

	"github.com/junioryono/Riot-API-Golang/constants/language"
	"github.com/junioryono/Riot-API-Golang/constants/patch"
)

type ProfileIcons []ProfileIcon

type ProfileIcon struct {
	ID    int   `json:"id"`
	Image Image `json:"image"`
}

func GetProfileIcons(v patch.Patch, lang language.Language) (ProfileIcons, error) {
	type Response struct {
		Data map[string]interface{} `json:"data"`
	}

	var res Response
	err := getJSON(fmt.Sprintf("http://ddragon.leagueoflegends.com/cdn/%s/data/%s/profileicon.json", v, lang), &res)

	var profileIcons ProfileIcons
	for _, item := range res.Data {
		profileIconMap := item.(map[string]interface{})
		idValue := profileIconMap["id"]

		var id int
		switch v := idValue.(type) {
		case string:
			id, err = strconv.Atoi(v)
			if err != nil {
				return profileIcons, err
			}
		case float64:
			id = int(v)
		default:
			return profileIcons, fmt.Errorf("unknown type for profileicon id: %T", v)
		}

		imageMap := profileIconMap["image"].(map[string]interface{})
		image := Image{
			Full:   imageMap["full"].(string),
			Group:  imageMap["group"].(string),
			Sprite: imageMap["sprite"].(string),
			H:      imageMap["h"].(float64),
			W:      imageMap["w"].(float64),
			Y:      imageMap["y"].(float64),
			X:      imageMap["x"].(float64),
		}

		profileIcon := ProfileIcon{
			ID:    id,
			Image: image,
		}

		profileIcons = append(profileIcons, profileIcon)
	}

	return profileIcons, err
}

func (pis ProfileIcons) ProfileIcon(profileIconID int) (ProfileIcon, error) {
	for _, pi := range pis {
		if pi.ID == profileIconID {
			return pi, nil
		}
	}

	return ProfileIcon{}, fmt.Errorf("profileicon %v not found", profileIconID)
}

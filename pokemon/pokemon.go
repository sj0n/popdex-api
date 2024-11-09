package pokemon

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"strings"

	"encore.dev/beta/errs"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type PokemonProfile struct {
	CacheControl string `header:"Cache-Control"`
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Weight       int    `json:"weight"`
	Height       int    `json:"height"`
	Abilities    []struct {
		Ability struct {
			Name string `json:"name"`
		} `json:"ability"`
	} `json:"abilities"`
	Sprites struct {
		FrontDefault string `json:"front_default"`
	} `json:"sprites"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

//encore:api public path=/pokemon/:name
func GetPokemonProfile(ctx context.Context, name string) (*PokemonProfile, error) {
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + name)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		switch resp.StatusCode {
		case 404:
			return nil, &errs.Error{
				Code:    errs.NotFound,
				Message: resp.Status,
			}
		case 400:
			return nil, &errs.Error{
				Code:    errs.InvalidArgument,
				Message: resp.Status,
			}
		default:
			return nil, &errs.Error{
				Code: errs.Internal,
			}
		}
	}

	var data PokemonProfile
	json.NewDecoder(resp.Body).Decode(&data)

	data.CacheControl = "public, max-age=604800, must-revalidate"
	return &data, nil
}

type GameVersionDetails struct {
	LevelLearnedAt  int `json:"level_learned_at"`
	MoveLearnMethod struct {
		Name string `json:"name"`
	} `json:"move_learn_method"`
	GameVersion struct {
		Name string `json:"name"`
	} `json:"version_group"`
}

type PokemonMoves struct {
	Moves []struct {
		Move struct {
			Name string `json:"name"`
		} `json:"move"`
		GameVersionDetails []GameVersionDetails `json:"version_group_details"`
	} `json:"moves"`
}

type FormattedMoves struct {
	Name        string `json:"name"`
	Level       int    `json:"level"`
	LearnMethod string `json:"learn_method"`
}

type GroupByVersion struct {
	CacheControl string                      `header:"Cache-Control"`
	Versions     map[string][]FormattedMoves `json:"versions"`
}

//encore:api public path=/pokemon/:name/moves
func GetPokemonMoves(ctx context.Context, name string) (*GroupByVersion, error) {
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + name)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		switch resp.StatusCode {
		case 404:
			return nil, &errs.Error{
				Code:    errs.NotFound,
				Message: resp.Status,
			}
		case 400:
			return nil, &errs.Error{
				Code:    errs.InvalidArgument,
				Message: resp.Status,
			}
		default:
			return nil, &errs.Error{
				Code: errs.Internal,
			}
		}
	}

	var data PokemonMoves
	json.NewDecoder(resp.Body).Decode(&data)

	groupedByVersion := make(map[string][]FormattedMoves)

	for i := 0; i < len(data.Moves); i++ {
		for j := 0; j < len(data.Moves[i].GameVersionDetails); j++ {
			groupedByVersion[data.Moves[i].GameVersionDetails[j].GameVersion.Name] = append(groupedByVersion[data.Moves[i].GameVersionDetails[j].GameVersion.Name], FormattedMoves{
				Name:        data.Moves[i].Move.Name,
				Level:       data.Moves[i].GameVersionDetails[j].LevelLearnedAt,
				LearnMethod: data.Moves[i].GameVersionDetails[j].MoveLearnMethod.Name,
			})
		}
	}

	return &GroupByVersion{Versions: groupedByVersion, CacheControl: "public, max-age=604800, must-revalidate"}, nil
}

type EncounterDetails struct {
	MinLevel int8 `json:"min_level"`
	MaxLevel int8 `json:"max_level"`
	Method   struct {
		Name string `json:"name"`
	} `json:"method"`
}

type VersionDetails struct {
	Version struct {
		Name string `json:"name"`
	} `json:"version"`
	EncounterDetails []EncounterDetails `json:"encounter_details"`
}

type Locations struct {
	LocationArea struct {
		Name string `json:"name"`
	} `json:"location_area"`
	VersionDetails []VersionDetails `json:"version_details"`
}

type FormattedLocations struct {
	LocationName string `json:"location_name"`
	MinLevel     int8   `json:"min_level"`
	MaxLevel     int8   `json:"max_level"`
	Method       string `json:"method"`
}

type PokemonLocations struct {
	CacheControl string                          `header:"Cache-Control"`
	Versions     map[string][]FormattedLocations `json:"versions"`
}

//encore:api public path=/pokemon/:name/locations
func GetPokemonLocations(ctx context.Context, name string) (*PokemonLocations, error) {
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + name + "/encounters")

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		switch resp.StatusCode {
		case 404:
			return nil, &errs.Error{
				Code:    errs.NotFound,
				Message: resp.Status,
			}
		case 400:
			return nil, &errs.Error{
				Code:    errs.InvalidArgument,
				Message: resp.Status,
			}
		default:
			return nil, &errs.Error{
				Code: errs.Internal,
			}
		}
	}

	var locations []Locations
	json.NewDecoder(resp.Body).Decode(&locations)

	type VersionLocation struct {
		Version  string
		Location string
		MinLevel int8
		MaxLevel int8
		Method   string
	}
	var flattenedData []VersionLocation

	for _, loc := range locations {
		for _, vDetail := range loc.VersionDetails {
			for _, encounter := range vDetail.EncounterDetails {
				titleFormatter := cases.Title(language.English)
				titledCase := titleFormatter.String(loc.LocationArea.Name)
				titledCase = strings.Join(strings.Split(titledCase, "-"), " ")

				flattenedData = append(flattenedData, VersionLocation{
					Version:  vDetail.Version.Name,
					Location: titledCase,
					MinLevel: encounter.MinLevel,
					MaxLevel: encounter.MaxLevel,
					Method:   encounter.Method.Name,
				})
			}
		}
	}

	result := make(map[string][]FormattedLocations)
	for i := 0; i < len(flattenedData); i++ {
		item := flattenedData[i]
		locations := result[item.Version]
		locationFound := false

		for j := 0; j < len(locations); j++ {
			existing := locations[j]
			if existing.LocationName == item.Location {
				locations[j].MinLevel = int8(math.Min(float64(existing.MinLevel), float64(item.MinLevel)))
				locations[j].MaxLevel = int8(math.Max(float64(existing.MaxLevel), float64(item.MaxLevel)))
				if !strings.Contains(existing.Method, item.Method) {
					locations[j].Method += ", " + item.Method
				}
				locationFound = true
				break
			}
		}

		if !locationFound {
			result[item.Version] = append(result[item.Version], FormattedLocations{
				LocationName: item.Location,
				MinLevel:     item.MinLevel,
				MaxLevel:     item.MaxLevel,
				Method:       item.Method,
			})
		}
	}

	return &PokemonLocations{Versions: result, CacheControl: "public, max-age=604800, must-revalidate"}, nil
}

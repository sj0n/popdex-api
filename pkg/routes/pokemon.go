package routes

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var titleFormatter = cases.Title(language.English)

type PokemonProfile struct {
	Id        int16  `json:"id"`
	Name      string `json:"name"`
	Weight    int16  `json:"weight"`
	Height    int16  `json:"height"`
	Abilities []struct {
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

func GetPokemonProfile(c echo.Context) error {
	name := c.Param("name")

	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + name)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		switch resp.StatusCode {
		case 404:
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Not Found",
			})
		case 400:
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Bad Request",
			})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Internal Server Error",
			})
		}
	}

	var result PokemonProfile

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal Server Error",
		})
	}

	c.Response().Header().Set("etag", resp.Header.Get("Etag"))
	fmt.Println("Route triggered")
	return c.JSON(http.StatusOK, result)
}

type GameVersionDetails struct {
	LevelLearnedAt  int8 `json:"level_learned_at"`
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
	Level       int8   `json:"level"`
	LearnMethod string `json:"learn_method"`
}

type GroupByVersion struct {
	Versions map[string][]FormattedMoves `json:"versions"`
}

func GetPokemonMoves(c echo.Context) error {
	name := c.Param("name")

	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal Server Error",
		})
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		switch resp.StatusCode {
		case 404:
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Not Found",
			})
		case 400:
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Bad Request",
			})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Internal Server Error",
			})
		}
	}

	var result PokemonMoves

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal Server Error",
		})
	}

	groupedByVersion := make(map[string][]FormattedMoves)

	for i := 0; i < len(result.Moves); i++ {
		for j := 0; j < len(result.Moves[i].GameVersionDetails); j++ {
			moveTitleCase := titleFormatter.String(result.Moves[i].Move.Name)
			moveTitleCase = strings.Join(strings.Split(moveTitleCase, "-"), " ")

			groupedByVersion[result.Moves[i].GameVersionDetails[j].GameVersion.Name] = append(groupedByVersion[result.Moves[i].GameVersionDetails[j].GameVersion.Name], FormattedMoves{
				Name:        moveTitleCase,
				Level:       result.Moves[i].GameVersionDetails[j].LevelLearnedAt,
				LearnMethod: result.Moves[i].GameVersionDetails[j].MoveLearnMethod.Name,
			})
		}
	}

	c.Response().Header().Set("etag", resp.Header.Get("Etag"))
	return c.JSON(http.StatusOK, &GroupByVersion{Versions: groupedByVersion})
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
	Versions map[string][]FormattedLocations `json:"versions"`
}

func GetPokemonLocations(c echo.Context) error {
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + c.Param("name") + "/encounters")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		switch resp.StatusCode {
		case 404:
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Not Found",
			})
		case 400:
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Bad Request",
			})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Internal Server Error",
			})
		}
	}

	var locations []Locations

	if err := json.NewDecoder(resp.Body).Decode(&locations); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal Server Error",
		})
	}

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

	c.Response().Header().Set("etag", resp.Header.Get("Etag"))
	return c.JSON(http.StatusOK, &PokemonLocations{Versions: result})
}

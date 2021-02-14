package web

import (
	"gomoney-mock-epl/database"
	"gomoney-mock-epl/teams"
	"net/http"

	"github.com/labstack/echo"
)

type CreateTeamRequest struct {
	City        string `json:"city"`
	HomeStadium string `json:"home_stadium"`
	LogoURL     string `json:"logo_url"`
	Name        string `json:"name"`
	NameAbbr    string `json:"name_abbr"`
	ShortName   string `json:"short_name"`
}

func createTeam(db teams.TeamsDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		dto := CreateTeamRequest{}
		if err := c.Bind(&dto); err != nil {
			return err
		}
		team, err := db.Create(c.Request().Context(), teams.Team{
			HomeStadium: dto.HomeStadium,
			LogoURL:     dto.LogoURL,
			Name:        dto.Name,
			NameAbbr:    dto.NameAbbr,
			ShortName:   dto.ShortName,
		})
		if err != nil {
			if database.IsDuplicateKeyError(err) {
				return echo.NewHTTPError(http.StatusConflict,
					errorDto("teams/already-exists", "This team already exists"))
			}
			return err
		}
		return c.JSON(http.StatusCreated,
			dataResponse("Team", "Team created successfully", team))
	}
}

func TeamRoutes(db teams.TeamsDB) RouteProvider {
	return func(e *echo.Echo) {
		e.POST("/teams/", createTeam(db), jwtMiddleware, onlyAdmins)
	}
}

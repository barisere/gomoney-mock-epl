package web

import (
	"fmt"
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

func (c *CreateTeamRequest) FromTeam(team teams.Team) *CreateTeamRequest {
	c.City = team.City
	c.HomeStadium = team.HomeStadium
	c.LogoURL = team.LogoURL
	c.Name = team.Name
	c.NameAbbr = team.NameAbbr
	c.ShortName = team.ShortName
	return c
}

func (c CreateTeamRequest) ToTeam(teamID string) teams.Team {
	return teams.Team{
		ID:          teamID,
		City:        c.City,
		HomeStadium: c.HomeStadium,
		LogoURL:     c.LogoURL,
		Name:        c.Name,
		NameAbbr:    c.NameAbbr,
		ShortName:   c.ShortName,
	}
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

func listTeams(db teams.TeamsDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		teams, err := db.List(c.Request().Context())
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK,
			dataResponse("Teams", "Available EPL teams", teams))
	}
}

func deleteTeam(db teams.TeamsDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		teamID := c.Param("team_id")
		if err := db.Delete(c.Request().Context(), teamID); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, nil)
	}
}

func editTeam(db teams.TeamsDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		teamID := c.Param("team_id")
		team, err := db.ByID(c.Request().Context(), teamID)
		if err != nil {
			return err
		}
		if team == nil {
			return echo.NewHTTPError(http.StatusNotFound,
				errorDto("NotFound", "That team does not exist"))
		}
		dto := (&CreateTeamRequest{}).FromTeam(*team)
		if err := c.Bind(dto); err != nil {
			return err
		}
		update := dto.ToTeam(team.ID)
		update.CreatedAt = team.CreatedAt
		team, err = db.Update(c.Request().Context(), update)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK,
			dataResponse("Team", "Team updated successfully", team))
	}
}

func viewTeam(db teams.TeamsDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		teamID := c.Param("team_id")
		team, err := db.ByID(c.Request().Context(), teamID)
		if err != nil {
			return err
		}
		if team == nil {
			return echo.NewHTTPError(http.StatusNotFound,
				errorDto("NotFound", "That team does not exist"))
		}
		return c.JSON(http.StatusOK,
			dataResponse("Team", fmt.Sprintf("Team: %q", team.Name), team))
	}
}
func TeamRoutes(db teams.TeamsDB) RouteProvider {
	return func(e *echo.Echo) {
		teams := e.Group("/teams", jwtMiddleware)
		teams.POST("/", createTeam(db), onlyAdmins)
		teams.GET("/", listTeams(db))
		teams.DELETE("/:team_id", deleteTeam(db), onlyAdmins)
		teams.GET("/:team_id", viewTeam(db))
		teams.PATCH("/:team_id", editTeam(db), onlyAdmins)
	}
}

package grifts

import (
	"encoding/json"
	"gomoney-mock-epl/config"
	"gomoney-mock-epl/database"
	"gomoney-mock-epl/fixtures"
	"gomoney-mock-epl/teams"
	"gomoney-mock-epl/users"
	"gomoney-mock-epl/web"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	. "github.com/markbates/grift/grift"
	"syreclabs.com/go/faker"
)

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

var _ = Namespace("db", func() {
	config, err := config.LoadConfig()
	panicOnErr(err)
	db, err := database.ConnectToDB(config.MongoURL)
	panicOnErr(err)
	app, err := web.NewApplication(db, *config)
	panicOnErr(err)

	Desc("reindex",
		"Recreate database indexes. This isn't a migration, drops and recreates the indexes")
	Add("reindex", func(c *Context) error {
		return database.CreateIndexes(db.Database(database.MockEPLDatabase))
	})

	Desc("create-admin", "Set up initial admin account")
	Add("create-admin", func(c *Context) error {
		_, err = users.SignUpAdmin(c, users.SignUpIntent{
			Email:     "superuser@gomoney-epl.local",
			FirstName: faker.Name().FirstName(),
			LastName:  faker.Name().LastName(),
			Password:  "password",
		}, app.AdminDB)

		return err
	})

	Desc("create-teams", "Seed database with teams")
	Add("create-teams", func(c *Context) error {
		seedTeams := Teams{}
		f, err := os.Open("grifts/teams.json")
		panicOnErr(err)
		defer f.Close()
		jsonBytes, err := ioutil.ReadAll(f)
		panicOnErr(err)
		err = json.Unmarshal(jsonBytes, &seedTeams)
		panicOnErr(err)
		tms := teamsFromSeedData(seedTeams)
		for _, t := range tms {
			if _, err := app.TeamsDB.Create(c, t); err != nil {
				return err
			}
		}
		return nil
	})

	Desc("create-fixtures", "Seed database with fixtures")
	Add("create-fixtures", func(c *Context) error {
		existingTeams, err := app.TeamsDB.List(c)
		panicOnErr(err)
		for i := 0; i < len(existingTeams)-1; i++ {
			req := makeFixtureRequest(existingTeams[i], existingTeams[i+1])
			if _, err := app.FixturesDB.Create(c, req); err != nil {
				return err
			}
		}
		return nil
	})

	Desc("fresh-setup", "Drop the existing database, recreate it, seed it with data")
	Add("fresh-setup", func(c *Context) error {
		if err := app.DefaultDB.Drop(c); err != nil {
			return err
		}
		return chainGrifts(c,
			"db:reindex",
			"db:create-admin",
			"db:create-teams",
			"db:create-fixtures")
	})
})

func chainGrifts(c *Context, grifts ...string) error {
	for _, g := range grifts {
		if err := Run(g, c); err != nil {
			return err
		}
	}
	return nil
}

func makeFixtureRequest(teamA teams.Team, teamB teams.Team) fixtures.CreateFixtureRequest {
	return fixtures.CreateFixtureRequest{
		HomeTeam:  teamA.ID,
		AwayTeam:  teamB.ID,
		MatchDate: time.Now().Add(time.Hour * time.Duration(rand.Intn(96))),
	}
}

func teamsFromSeedData(ts Teams) []teams.Team {
	var tms = make([]teams.Team, 0, len(ts))
	for _, t := range ts {
		tms = append(tms, toDomainTeam(t))
	}
	return tms
}

func toDomainTeam(t Team) teams.Team {
	return teams.Team{
		City:        t.Grounds[0].City,
		HomeStadium: t.Grounds[0].Name,
		Name:        t.Club.Name,
		NameAbbr:    t.Club.Abbr,
		ShortName:   t.Club.ShortName,
	}
}

package grifts

// Generated by https://quicktype.io

type Teams []Team

type Team struct {
	Name      string   `json:"name"`
	Club      Club     `json:"club"`
	TeamType  TeamType `json:"teamType"`
	Grounds   []Ground `json:"grounds"`
	ShortName string   `json:"shortName"`
	ID        float64  `json:"id"`
	AltIDS    AltIDS   `json:"altIds"`
}

type AltIDS struct {
	Opta string `json:"opta"`
}

type Club struct {
	Name      string  `json:"name"`
	ShortName string  `json:"shortName"`
	Abbr      string  `json:"abbr"`
	ID        float64 `json:"id"`
}

type Ground struct {
	Name     string    `json:"name"`
	City     string    `json:"city"`
	Capacity *float64  `json:"capacity,omitempty"`
	Source   *Source   `json:"source,omitempty"`
	ID       float64   `json:"id"`
	Location *Location `json:"location,omitempty"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Source string

const (
	Opta Source = "OPTA"
)

type TeamType string

const (
	First TeamType = "FIRST"
)
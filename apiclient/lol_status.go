package apiclient

import (
	"time"

	"github.com/junioryono/Riot-API-Golang/apiclient/ratelimiter"
	"github.com/junioryono/Riot-API-Golang/constants/region"
)

type StatusPlatformData struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Locales      []string      `json:"locales"`
	Maintenances []Maintenance `json:"maintenances"`
	Incidents    []Incident    `json:"incidents"`
}

type Maintenance struct {
	MaintenanceStatus string     `json:"maintenance_status"`
	Titles            []Title    `json:"titles"`
	ArchiveAt         *time.Time `json:"archive_at"`
	Updates           []Update   `json:"updates"`
	IncidentSeverity  *string    `json:"incident_severity"`
	UpdatedAt         *time.Time `json:"updated_at"`
	Platforms         []string   `json:"platforms"`
	ID                int        `json:"id"`
	CreatedAt         time.Time  `json:"created_at"`
}

type Incident struct {
	MaintenanceStatus string   `json:"maintenance_status"`
	CreatedAt         string   `json:"created_at"` // Note: time format is not clear from example
	UpdatedAt         *string  `json:"updated_at"` // Note: time format is not clear from example
	ID                int      `json:"id"`
	Titles            []Title  `json:"titles"`
	Updates           []Update `json:"updates"`
	Platforms         []string `json:"platforms"`
	IncidentSeverity  string   `json:"incident_severity"`
	ArchiveAt         *string  `json:"archive_at"` // Note: time format is not clear from example
}

type Title struct {
	Locale  string `json:"locale"`
	Content string `json:"content"`
}

type Update struct {
	Author           string   `json:"author"`
	PublishLocations []string `json:"publish_locations"`
	UpdatedAt        string   `json:"updated_at"` // Note: time format is not clear from example
	Publish          bool     `json:"publish"`
	ID               int      `json:"id"`
	Translations     []Title  `json:"translations"`
	CreatedAt        string   `json:"created_at"` // Note: time format is not clear from example
}

func (c *uniqueClient) GetStatusPlatformData(r region.Region) (*StatusPlatformData, error) {
	var res StatusPlatformData
	err := c.dispatchAndUnmarshal(r, "/lol/status/v4/platform-data", "", nil, ratelimiter.GetStatusPlatformData, &res)
	return &res, err
}

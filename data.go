package postmortem

type Geography interface {
	GeoID() string
	StateID() string
	RecordNo() int
}

// State represents an individual United States state.
type State struct {
	ID       string    `json:"id"`
	Abbrev   string    `json:"state"`
	Name     string    `json:"name"`
	RecNo    int       `json:"record_number"`
	Counties []*County `json:"counties"`
}

func (s *State) GeoID() string   { return s.ID }
func (s *State) StateID() string { return s.Abbrev }
func (s *State) RecordNo() int   { return s.RecNo }

// County represents an individual United States county.
type County struct {
	ID    string `json:"id"`
	State string `json:"state"`
	Name  string `json:"name"`
	RecNo int    `json:"record_number"`
}

func (c *County) GeoID() string   { return c.ID }
func (c *County) StateID() string { return c.State }
func (c *County) RecordNo() int   { return c.RecNo }

// ACSStatistics aggregates various statistics about a geography
// collected from the American Community Survey.
type ACSStatistics struct {
	TotalPopulation *TotalPopulation `json:"total_population"`
	FoodStamps      *FoodStamps      `json:"food_stamps"`
}

// TotalPopulation provides an estimate of the total population
// living within a geography.
type TotalPopulation struct {
	Total int `json:"total"`
}

// FoodStamps describes statistics about households receiving
// food stamps within the past 12 months.
type FoodStamps struct {
	Total int `json:"-"`
	Yes   int `json:"yes"`
	No    int `json:"no"`
}

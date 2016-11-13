package postmortem

// State represents an individual United States state.
type State struct {
	ID       string   `json:"id"`
	Abbrev   string   `json:"state"`
	Name     string   `json:"name"`
	RecNo    int      `json:"record_number"`
	Counties []County `json:"counties"`
}

// County represents an individual United States county.
type County struct {
	ID    string `json:"id"`
	State string `json:"state"`
	Name  string `json:"name"`
	RecNo int    `json:"record_number"`
}

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

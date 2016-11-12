package postmortem

// County represents an individual United States county.
type County struct {
	ID    string `json:"id"`
	State string `json:"state"`
	Name  string `json:"name"`
	RecNo int    `json:"record_number"` // lolwat
}

// FoodStamps describes statistics about households receiving
// food stamps within the past 12 months.
type FoodStamps struct {
	Total int `json:"-"`
	Yes   int `json:"yes"`
	No    int `json:"no"`
}

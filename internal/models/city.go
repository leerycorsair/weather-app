package models

// City represents a city model
// @Description City model
type City struct {
	Id        int     `json:"id"  db:"id"`               // @Description City ID
	Name      string  `json:"name"  db:"name"`           // @Description City name
	Country   string  `json:"country"  db:"country"`     // @Description Country name
	Latitude  float64 `json:"latitude"  db:"latitude"`   // @Description Latitude of the city
	Longitude float64 `json:"longitude"  db:"longitude"` // @Description Longitude of the city
}

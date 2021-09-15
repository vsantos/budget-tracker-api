package services

const (
	// DatabaseURI will set where to store data
	DatabaseURI = "mongodb://localhost:27017/"
)

// DataManager will define Databae operations
type DataManager interface {
	Health() (err error)
}

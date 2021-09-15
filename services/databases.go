package services

// DatabaseManager will define instance operations
type DatabaseManager interface {
	Health() (err error)
}

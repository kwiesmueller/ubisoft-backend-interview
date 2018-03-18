package feedback

// Repository interface for storing feedback
type Repository interface {
	Add(Entry) error
	GetLatest(n uint) ([]Entry, error)
	GetLatestFiltered(n uint, filter int) ([]Entry, error)
}

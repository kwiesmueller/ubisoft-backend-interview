package feedback

import (
	"github.com/playnet-public/libs/log"
	"go.uber.org/zap"
)

// Service for getting feedback
type Service struct {
	*log.Logger
	repo Repository
}

// New Service for getting feedback
func New(log *log.Logger, repo Repository) *Service {
	log = log.WithFields(zap.String("component", "feedback.service"))
	return &Service{
		Logger: log,
		repo:   repo,
	}
}

// Add entry to Repository
func (s *Service) Add(entry Entry) error {
	if entry.Rating > 5 || entry.Rating < 1 {
		return ErrInvalidRating
	}
	return s.repo.Add(entry)
}

// GetLatest n entries from Repository
func (s *Service) GetLatest(n uint) ([]Entry, error) {
	return s.repo.GetLatest(n)
}

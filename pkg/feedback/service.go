package feedback

import (
	"github.com/playnet-public/libs/log"
	"go.uber.org/zap"
)

// Service for getting feedback
type Service struct {
	*log.Logger
}

// New Service for getting feedback
func New(log *log.Logger) *Service {
	log = log.WithFields(zap.String("component", "feedback.service"))
	return &Service{
		Logger: log,
	}
}

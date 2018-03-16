package feedback

import (
	"testing"

	"github.com/playnet-public/libs/log"
)

func TestNew(t *testing.T) {
	log := log.New("nop", "", true)
	if New(log) == nil {
		t.Errorf("New() == nil")
	}
}

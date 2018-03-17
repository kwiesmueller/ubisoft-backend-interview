package feedback

import "errors"

var (
	// ErrInvalidRating .
	ErrInvalidRating = errors.New("rating invalid. has to be between 1-5")
	// ErrDuplicateEntry .
	ErrDuplicateEntry = errors.New("entries may only be supplied once per user/session")
)

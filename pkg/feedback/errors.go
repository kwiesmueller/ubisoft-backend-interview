package feedback

import "errors"

var (
	// ErrInvalidRating .
	ErrInvalidRating = errors.New("rating invalid. has to be between 1-5")
	// ErrDuplicateEntry .
	ErrDuplicateEntry = errors.New("entries may only be sent once per user/session")
	// ErrNoSession .
	ErrNoSession = errors.New("no sessionID provided")
	// ErrNoUserID .
	ErrNoUserID = errors.New("no userID provided")
)

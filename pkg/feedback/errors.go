package feedback

import "errors"

var (
	errInvalidRating = errors.New("rating invalid. has to be between 1-5")
)

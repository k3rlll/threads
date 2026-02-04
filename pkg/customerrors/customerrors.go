package customerrors

import "errors"

var (
	ErrNoTagsAffected = errors.New("no rows were affected by the operation")
)

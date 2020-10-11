package storage

import "errors"

var (
	ErrUnallowedField = errors.New("storage: unallowed field")
	ErrNoRowsAffected = errors.New("storage: no rows affected")
)

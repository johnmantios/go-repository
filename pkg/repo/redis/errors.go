package redis

import "errors"

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrEditConflict    = errors.New("edit conflict")
	ErrChildlessRecord = errors.New("record has no related children")
)

package fs

import (
	"errors"
	"os"
)

var (
	ErrIsDir            = errors.New("is a directory")
	ErrNotDir           = errors.New("not a directory")
	ErrIsFile           = errors.New("is a file")
	ErrNotFile          = errors.New("not a file")
	ErrInvalid          = os.ErrInvalid
	ErrPermission       = os.ErrPermission
	ErrExist            = os.ErrExist
	ErrNotExist         = os.ErrNotExist
	ErrClosed           = os.ErrClosed
	ErrNoDeadline       = os.ErrNoDeadline
	ErrDeadlineExceeded = os.ErrDeadlineExceeded
)

var (
	PathSeparator = string(os.PathSeparator)
)

// ----------------------------------------------------------------------------
// UTILS
// ----------------------------------------------------------------------------

// func Or[T any](value T, err error, default_ T) T {
// 	if err != nil {
// 		return default_
// 	}
// 	return value
// }

func Force[T any](value T, err error) T {
	return value
}

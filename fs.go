// Package fs provides utilities for file system operations, including error handling,
// directory management, path manipulations, and event monitoring.
package fs

import (
	"errors"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
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

type Event struct {
	Op   fsnotify.Op
	Path string
	Err  error
}

func (e Event) Has(op fsnotify.Op) bool {
	return e.Op.Has(op)
}

func (e Event) String() string {
	res := e.Op.String()
	if e.Has(EvtError) {
		res += "|Error"
	}
	return strings.TrimPrefix(res, "|")
}

var (
	EvtCreate = fsnotify.Create
	EvtRemove = fsnotify.Remove
	EvtWrite  = fsnotify.Write
	EvtRename = fsnotify.Rename
	EvtChmod  = fsnotify.Chmod
	EvtError  = fsnotify.Op(2048)
)

var (
	PathSeparator = string(os.PathSeparator)
)

func Force[T any](value T, err error) T {
	return value
}

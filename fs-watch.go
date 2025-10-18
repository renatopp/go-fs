package fs

import (
	"context"
	"slices"

	"github.com/fsnotify/fsnotify"
)

// Watcher represents a file system watcher.
type Watcher struct {
	watcher *fsnotify.Watcher
	files   []string
}

// NewWatcher creates a new file system watcher.
func NewWatcher() (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		watcher: w,
		files:   []string{},
	}, nil
}

// Watch starts watching for file system events and invokes the provided
// callback function for each event.
func (w *Watcher) Watch(ctx context.Context, callback func(event Event)) error {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return nil
			}
			callback(Event{
				Op:   event.Op,
				Path: event.Name,
			})
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return err
			}
			callback(Event{
				Op:  EvtError,
				Err: err,
			})
		case <-ctx.Done():
			return nil
		}
	}
}

// Add adds a path to the watcher.
func (w *Watcher) Add(p string) error {
	err := w.watcher.Add(p)
	if err == nil {
		w.files = append(w.files, p)
	}
	return err
}

// Has checks if a path is being watched.
func (w *Watcher) Has(path string) bool {
	return slices.Contains(w.files, path)
}

// Remove removes a path from the watcher.
func (w *Watcher) Remove(p string) error {
	err := w.watcher.Remove(p)
	if err == nil {
		for i, file := range w.files {
			if file == p {
				w.files = append(w.files[:i], w.files[i+1:]...)
				break
			}
		}
	}
	return err
}

// WatchList returns the list of paths being watched.
func (w *Watcher) WatchList() []string {
	return w.watcher.WatchList()
}

// Close closes the watcher.
func (w *Watcher) Close() error {
	return w.watcher.Close()
}

// Watch watches a path for file system events and invokes the provided
// callback function for each event.
func Watch(ctx context.Context, p string, callback func(event Event)) error {
	w, err := NewWatcher()
	if err != nil {
		return err
	}
	w.Add(p)
	return w.Watch(ctx, callback)
}

// WatchRecursive watches a path and all its subdirectories for file system
// events and invokes the provided callback function for each event.
func WatchRecursive(ctx context.Context, p string, callback func(event Event)) error {
	w, err := NewWatcher()
	if err != nil {
		return err
	}
	w.Add(p)
	return w.Watch(ctx, func(event Event) {
		if event.Has(EvtCreate) && IsDir(event.Path) {
			w.Add(event.Path)
		}
		if event.Has(EvtRemove) && IsDir(event.Path) {
			w.Remove(event.Path)
		}
		if event.Has(EvtRename) && IsDir(event.Path) {
			w.Remove(event.Path)
		}
		callback(event)
	})
}

// WatchGlob watches a directory for file system events matching a glob pattern
// and invokes the provided callback function for each matching event.
func WatchGlob(ctx context.Context, dir string, pattern string, callback func(event Event)) error {
	if !IsPatternValid(pattern) {
		return ErrInvalid
	}

	return WatchRecursive(ctx, dir, func(event Event) {
		if ForceMatch(ToSlashPath(event.Path), JoinPathLinux(dir, pattern)) {
			callback(event)
		}
	})
}

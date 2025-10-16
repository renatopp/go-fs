package fs

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
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

// ----------------------------------------------------------------------------
// DIRECTORY-RELATED OPERATIONS
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// GENERAL OPERATIONS (FILES AND DIRECTORIES)
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// PATH OPERATIONS
// ----------------------------------------------------------------------------

func PathSeparator() string {
	return string(os.PathSeparator)
}

func JoinPath(elem ...string) string {
	return filepath.Join(elem...)
}

func JoinPathLinux(elem ...string) string {
	return strings.Join(elem, "/")
}

func JoinPathWindows(elem ...string) string {
	return strings.Join(elem, "\\")
}

func JoinPathWith(sep string, elem ...string) string {
	return strings.Join(elem, sep)
}

func AbsolutePath(p string) (string, error) {
	absPath, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	return absPath, nil
}

func RelativePath(base, target string) (string, error) {
	absBase, err := AbsolutePath(base)
	if err != nil {
		return "", err
	}
	absTarget, err := AbsolutePath(target)
	if err != nil {
		return "", err
	}
	relPath, err := filepath.Rel(absBase, absTarget)
	if err != nil {
		return "", err
	}
	return relPath, nil
}

func ExtensionPath(p string) string {
	return filepath.Ext(p)
}

func HasExtensionPath(p string) bool {
	return filepath.Ext(p) != ""
}

func StemPath(p string) string {
	ext := filepath.Ext(p)
	return strings.TrimSuffix(filepath.Base(p), ext)
}

func BasePath(p string) string {
	return filepath.Base(p)
}

func DirPath(p string) string {
	return filepath.Dir(p)
}

func VolumePath(p string) string {
	return filepath.VolumeName(p)
}

func SplitPath(p string) (dir, file string) {
	return filepath.Split(p)
}

func SplitExtPath(p string) (string, string) {
	ext := filepath.Ext(p)
	stem := strings.TrimSuffix(filepath.Base(p), ext)
	dir := filepath.Dir(p)
	return filepath.Join(dir, stem), ext
}

func ParentPathName(p string) string {
	dir := filepath.Dir(p)
	return filepath.Base(dir)
}

func IsAbsolutePath(p string) bool {
	return filepath.IsAbs(p)
}

func ToSlashPath(p string) string {
	return filepath.ToSlash(p)
}

func CleanPath(p string) string {
	return filepath.Clean(p)
}

func GetDir(p string) (string, error) {
	if !Exists(p) {
		return "", os.ErrNotExist
	}

	if IsFile(p) {
		return filepath.Dir(p), nil
	}

	return p, nil
}

func GetParentDirName(p string) (string, error) {
	dir, err := GetDir(p)
	if err != nil {
		return "", err
	}
	parent := filepath.Dir(dir)
	if parent == dir {
		return "", nil
	}
	return filepath.Base(parent), nil
}

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

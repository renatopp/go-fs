package fs

import (
	"path/filepath"
	"strings"
)

// PathParts represents the different components of a file path.
type PathParts struct {
	Absolute   string // Absolute path (eg: /home/users/dev/fs/path.go)
	Base       string // Base name (eg: path.go)
	Name       string // Name without extension (eg: path)
	Ext        string // Extension with dot (eg: .go)
	ExtName    string // Extension without dot (eg: go)
	Parent     string // Parent directory (eg: /home/users/dev/fs)
	ParentName string // Parent directory name (eg: fs)
	Volume     string // Volume name (eg: C: on Windows, empty on Unix)
}

// JoinPath joins any number of path elements into a single path,
// using the OS-specific path separator.
func JoinPath(elem ...string) string {
	return filepath.Join(elem...)
}

// JoinPathUnix joins any number of path elements into a single path,
// using the Unix-style path separator (/).
func JoinPathLinux(elem ...string) string {
	return strings.Join(elem, "/")
}

// JoinPathWindows joins any number of path elements into a single path,
// using the Windows-style path separator (\).
func JoinPathWindows(elem ...string) string {
	return strings.Join(elem, "\\")
}

// JoinPathURL joins any number of path elements into a single path,
// using the provided separator.
func JoinPathWith(sep string, elem ...string) string {
	return strings.Join(elem, sep)
}

// AbsolutePath returns the absolute path for the given path.
func AbsolutePath(p string) (string, error) {
	absPath, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	return absPath, nil
}

// ForceAbsolutePath is like AbsolutePath but ignores any errors and returns
// the absolute path.
func ForceAbsolutePath(p string) string {
	abs, _ := AbsolutePath(p)
	return abs
}

// CleanAbsolutePath returns the cleaned absolute path for the given path.
func RelativePath(base, target string) (string, error) {
	absBase, err := AbsolutePath(base)
	if err != nil {
		return target, err
	}
	absTarget, err := AbsolutePath(target)
	if err != nil {
		return target, err
	}
	relPath, err := filepath.Rel(absBase, absTarget)
	if err != nil {
		return target, err
	}
	return relPath, nil
}

// ForceRelativePath is like RelativePath but ignores any errors and returns
// the relative path.
func ForceRelativePath(base, target string) string {
	rel, _ := RelativePath(base, target)
	return rel
}

// IsAbsolutePath checks if the given path is an absolute path.
func IsAbsolutePath(p string) bool {
	return filepath.IsAbs(p)
}

// CleanPath cleans the given path, resolving any . or .. elements and removing
// any redundant separators.
func CleanPath(p string) string {
	return filepath.Clean(p)
}

// ToBackslashPath converts all forward slashes in the path to backslashes.
func ToBackslashPath(p string) string {
	return strings.ReplaceAll(p, "/", "\\")
}

// ToSlashPath converts all backslashes in the path to forward slashes.
func ToSlashPath(p string) string {
	return filepath.ToSlash(p)
}

// IsSlashPath checks if the given path contains any forward slashes.
func IsSlashPath(p string) bool {
	return strings.Contains(p, "/")
}

// IsBackslashPath checks if the given path contains any backslashes.
func IsBackslashPath(p string) bool {
	return strings.Contains(p, "\\")
}

// HasExtensionPath checks if the given path has a file extension.
func HasExtensionPath(p string) bool {
	return filepath.Ext(p) != ""
}

// SplitPath splits the given path into its components using forward slashes
// as the separator. It returns a slice of path components. If you need to
// get the path parts, use GetPathParts instead.
func SplitPath(p string) []string {
	return strings.Split(ToSlashPath(p), "/")
}

// Returns the base of the file or directory specified by path. Example:
// for path "/home/user/file.txt", it returns "file.txt".
func GetPathBase(p string) string {
	return filepath.Base(p)
}

// Returns the name of the file or directory specified by path without the
// extension. Example: for path "/home/user/file.txt", it returns "file".
func GetPathName(p string) string {
	ext := filepath.Ext(p)
	return strings.TrimSuffix(filepath.Base(p), ext)
}

// Returns the extension of the file specified by path. Example: for path
// "/home/user/file.txt", it returns ".txt".
func GetPathExtension(p string) string {
	return filepath.Ext(p)
}

// Returns the extension name of the file specified by path without the dot.
// Example: for path "/home/user/file.txt", it returns "txt".
func GetPathExtensionName(p string) string {
	ext := filepath.Ext(p)
	return strings.TrimPrefix(ext, ".")
}

// Returns the parent directory of the file or directory specified by path.
// Example: for path "/home/user/file.txt", it returns "/home/user".
func GetPathParent(p string) string {
	return filepath.Dir(p)
}

// Returns the name of the parent directory of the file or directory specified
// by path. Example: for path "/home/user/file.txt", it returns "user".
func GetPathParentName(p string) string {
	dir := filepath.Dir(p)
	return filepath.Base(dir)
}

// Returns the volume name of the file or directory specified by path.
// Example: for path "C:\Users\user\file.txt", it returns "C:" on Windows,
// and for path "/home/user/file.txt", it returns an empty string on Unix.
func GetPathVolume(p string) string {
	return filepath.VolumeName(p)
}

// GetPathParts returns the different components of a file path as a PathParts struct.
func GetPathParts(p string) PathParts {
	abs := ForceAbsolutePath(p)
	return PathParts{
		Absolute:   abs,
		Base:       GetPathBase(abs),
		Name:       GetPathName(abs),
		Ext:        GetPathExtension(abs),
		ExtName:    GetPathExtensionName(abs),
		Parent:     GetPathParent(abs),
		ParentName: GetPathParentName(abs),
		Volume:     GetPathVolume(abs),
	}
}

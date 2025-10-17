package fs

import (
	"path/filepath"
	"strings"
)

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

func ForceAbsolutePath(p string) string {
	abs, _ := AbsolutePath(p)
	return abs
}

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

func ForceRelativePath(base, target string) string {
	rel, _ := RelativePath(base, target)
	return rel
}

func IsAbsolutePath(p string) bool {
	return filepath.IsAbs(p)
}

func CleanPath(p string) string {
	return filepath.Clean(p)
}

func ToBackslashPath(p string) string {
	return strings.ReplaceAll(p, "/", "\\")
}

func FromBackslashPath(p string) string {
	return strings.ReplaceAll(p, "\\", "/")
}

func ToSlashPath(p string) string {
	return filepath.ToSlash(p)
}

func FromSlashPath(p string) string {
	return filepath.FromSlash(p)
}

func IsSlashPath(p string) bool {
	return strings.Contains(p, "/")
}

func IsBackslashPath(p string) bool {
	return strings.Contains(p, "\\")
}

func HasExtensionPath(p string) bool {
	return filepath.Ext(p) != ""
}

func SplitPath(p string) []string {
	return strings.Split(ToSlashPath(p), "/")
}

func GetPathBase(p string) string {
	return filepath.Base(p)
}

func GetPathName(p string) string {
	ext := filepath.Ext(p)
	return strings.TrimSuffix(filepath.Base(p), ext)
}

func GetPathExtension(p string) string {
	return filepath.Ext(p)
}

func GetPathExtensionName(p string) string {
	ext := filepath.Ext(p)
	return strings.TrimPrefix(ext, ".")
}

func GetPathParent(p string) string {
	return filepath.Dir(p)
}

func GetPathParentName(p string) string {
	dir := filepath.Dir(p)
	return filepath.Base(dir)
}

func GetPathVolume(p string) string {
	return filepath.VolumeName(p)
}

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

package fs

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"hash"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/bmatcuk/doublestar/v4"
)

// Exists checks if a file or directory exists at the given p.
func Exists(p string) bool {
	_, err := os.Stat(p)
	return !os.IsNotExist(err)
}

// IsEmpty checks if a file or directory at the specified path is empty.
// For files, it checks if the file size is zero bytes.
// For directories, it checks if the directory contains no files or subdirectories.
// If the path does not exist, it returns an error but also true.
func IsEmpty(p string) (bool, error) {
	if IsDir(p) {
		return isEmptyDir(p)
	} else if IsFile(p) {
		return isEmptyFile(p)
	} else {
		return true, os.ErrNotExist
	}
}

// ForceIsEmpty is like IsEmpty but ignores any errors and returns false in such cases.
func ForceIsEmpty(p string) bool {
	empty, _ := IsEmpty(p)
	return empty
}

// IsFile checks if two files or directories at the specified paths refer to the same file
// or directory.
func IsSame(p1, p2 string) bool {
	s1, err := os.Stat(p1)
	if err != nil {
		return false
	}
	s2, err := os.Stat(p2)
	if err != nil {
		return false
	}
	return os.SameFile(s1, s2)
}

// IsExecutable checks if a file at the specified path is executable. Directories
// are considered not executable.
func IsExecutable(p string) bool {
	if !IsFile(p) {
		return false
	}
	info, err := os.Stat(p)
	if err != nil {
		return false
	}
	mode := info.Mode()
	if mode&0111 != 0 {
		return true
	}
	if runtime.GOOS == "windows" {
		ext := strings.ToLower(filepath.Ext(p))
		pathext := os.Getenv("PATHEXT")
		for _, e := range strings.Split(pathext, ";") {
			if strings.ToLower(e) == ext {
				return true
			}
		}
	}
	return false
}

// IsReadable checks if a file at the specified path is readable. Directories
// are considered not readable.
func IsReadable(p string) bool {
	if !IsFile(p) {
		return false
	}
	file, err := os.OpenFile(p, os.O_RDONLY, 0)
	if err != nil {
		return false
	}
	file.Close()
	return true
}

// IsWritable checks if a file at the specified path is writable. Directories
// are considered not writable.
func IsWritable(p string) bool {
	if !IsFile(p) {
		return false
	}
	file, err := os.OpenFile(p, os.O_WRONLY, 0)
	if err != nil {
		return false
	}
	file.Close()
	return true
}

// IsHidden checks if a file or directory at the specified path is hidden.
//
// On Unix-like systems, a file or directory is considered hidden if its name
// starts with a dot ('.').
//
// On Windows, a file or directory is considered hidden if it has the
// FILE_ATTRIBUTE_HIDDEN attribute set.
func IsHidden(p string) (bool, error) {
	abs := Force(AbsolutePath(p))
	base := filepath.Base(p)
	if runtime.GOOS == "windows" {
		pointer, err := syscall.UTF16PtrFromString(abs)
		if err != nil {
			return false, err
		}
		attributes, err := syscall.GetFileAttributes(pointer)
		if err != nil {
			return false, err
		}
		return attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0, nil
	}
	return strings.HasPrefix(base, "."), nil
}

// ForceIsHidden is like IsHidden but ignores any errors and returns false in
// such cases.
func ForceIsHidden(p string) bool {
	hidden, _ := IsHidden(p)
	return hidden
}

// IsPatternValid checks if a glob pattern is valid.
func IsPatternValid(pattern string) bool {
	return doublestar.ValidatePattern(pattern)
}

// Walk traverses the directory tree rooted at the specified path, calling the
// provided function for each file or directory encountered. The function receives
// the relative path of the files or directories found as its argument. If the callback
// function returns an error, the walk is aborted and the error is returned.
func Walk(p string, fn func(string) error) error {
	return filepath.WalkDir(p, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		return fn(ForceRelativePath(p, path))
	})
}

// List returns a slice of names of all entries (files and directories) within
// the specified directory path. If the directory does not exist or is not
// accessible, it returns an error. This function does not include the full paths,
// only the names of the entries.
//
// This function is not recursive; it only lists entries in the specified
// directory, not in its subdirectories.
func List(p string) ([]string, error) {
	entries, err := os.ReadDir(p)
	files := []string{}
	if err != nil {
		return files, err
	}
	names := make([]string, len(entries))
	for i, entry := range entries {
		names[i] = entry.Name()
	}
	return names, nil
}

// ForceList is like List but ignores any errors and returns an empty slice
// in such cases.
func ForceList(p string) []string {
	list, _ := List(p)
	return list
}

// ListRecursive returns a slice of relative paths of all entries (files and
// directories) within the specified directory path and its subdirectories.
// If the directory does not exist or is not accessible, it returns an error.
// The returned paths are relative to the specified directory.
//
// This function is recursive; it lists entries in the specified directory
// and all its subdirectories.
func ListRecursive(p string) ([]string, error) {
	if !IsDir(p) {
		return nil, ErrNotDir
	}

	results := []string{}
	err := filepath.WalkDir(p, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(p, path)
		if err != nil {
			return err
		}
		if relPath != "." {
			results = append(results, relPath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// ForceListRecursive is like ListRecursive but ignores any errors and returns
// an empty slice in such cases.
func ForceListRecursive(p string) []string {
	list, _ := ListRecursive(p)
	return list
}

// Glob returns the names of all files matching pattern or nil if there is no
// matching file. The syntax of patterns is the same as in filepath.Match.
// The pattern may describe hierarchical names such as /usr/*/bin/ed (assuming
// the Separator is '/').
func Glob(dir, pattern string) ([]string, error) {
	files, err := doublestar.Glob(nil, filepath.Join(dir, pattern))
	if files == nil {
		files = []string{}
	}
	r := len(dir) + len(PathSeparator)
	for i, f := range files {
		files[i] = f[r:]
	}
	return files, err
}

// ForceGlob is like Glob but ignores any errors and returns an empty slice in
// such cases.
func ForceGlob(dir string, pattern string) []string {
	files, _ := Glob(dir, pattern)
	return files
}

// Match reports whether the path matches the given pattern. The pattern syntax
// is the [doublestar](https://github.com/bmatcuk/doublestar#readme) syntax.
func Match(p, pattern string) (bool, error) {
	return doublestar.Match(pattern, p)
}

// ForceMatch is like Match but ignores any errors and returns false in such
// cases.
func ForceMatch(p, pattern string) bool {
	matched, _ := Match(p, pattern)
	return matched
}

// Copy copies a file or directory from src to dst. If src is a directory, it
// copies the entire directory recursively. If src is a file, it copies the file.
// If dst does not exist, it will be created. If it exists, it will be merged
// (for directories) or overwritten (for files).
func Copy(src, dst string) error {
	if IsDir(src) {
		return copyDir(src, dst)
	}
	return copyFile(src, dst)
}

// Move moves a file or directory from src to dst. It is equivalent to renaming
// the file or directory. If src and dst are on different filesystems, it
// performs a copy followed by a delete of the original.
func Move(src, dst string) error {
	return os.Rename(src, dst)
}

// Rename renames (moves) a file or directory from oldPath to newPath. If oldPath
// and newPath are on different filesystems, it performs a copy followed by a
// delete of the original.
func Rename(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

// Remove removes a file or directory at the specified path. If the path is a
// directory, it removes the directory and all its contents recursively.
// If the path does not exist, it returns nil (no error).
// If there is an error, it will be of type [*PathError].
func Remove(p string) error {
	return os.RemoveAll(p)
}

// SetMode sets the file mode (permissions) of a file at the specified path. If
// the path does not exist, it returns an error.
func SetMode(p string, mode os.FileMode) error {
	return os.Chmod(p, mode)
}

// SetHidden sets or unsets the hidden attribute of a file or directory at the
// specified path.
//
// On Unix-like systems, it renames the file or directory to start with a dot ('.')
// to hide it, or removes the leading dot to unhide it.
//
// On Windows, it sets or clears the FILE_ATTRIBUTE_HIDDEN attribute.
func SetHidden(p string, hidden bool) error {
	abs := Force(AbsolutePath(p))
	if runtime.GOOS == "windows" {
		pointer, err := syscall.UTF16PtrFromString(abs)
		if err != nil {
			return err
		}
		attributes, err := syscall.GetFileAttributes(pointer)
		if err != nil {
			return err
		}
		if hidden {
			attributes |= syscall.FILE_ATTRIBUTE_HIDDEN
		} else {
			attributes &^= syscall.FILE_ATTRIBUTE_HIDDEN
		}
		return syscall.SetFileAttributes(pointer, attributes)
	}
	base := filepath.Base(p)
	dir := filepath.Dir(p)
	if hidden {
		if strings.HasPrefix(base, ".") {
			return nil
		}
		newPath := filepath.Join(dir, "."+base)
		return os.Rename(p, newPath)
	} else {
		if !strings.HasPrefix(base, ".") {
			return nil
		}
		newBase := strings.TrimPrefix(base, ".")
		newPath := filepath.Join(dir, newBase)
		return os.Rename(p, newPath)
	}
}

// Hide sets the hidden attribute of a file or directory at the specified path.
// Same as SetHidden with hidden=true.
func Hide(p string) error {
	return SetHidden(p, true)
}

// Unhide clears the hidden attribute of a file or directory at the specified path.
// Same as SetHidden with hidden=false.
func Unhide(p string) error {
	return SetHidden(p, false)
}

// Chmod is an alias for SetMode.
func Chmod(p string, mode os.FileMode) error {
	return os.Chmod(p, mode)
}

// Chown changes the ownership of a file at the specified path to the given
// user ID (uid) and group ID (gid). If the path does not exist, it returns
// an error.
func Chown(p string, uid, gid int) error {
	return os.Chown(p, uid, gid)
}

// Chdir changes the current working directory to the specified path. If the
// path does not exist or is not a directory, it returns an error.
func Chdir(p string) error {
	return os.Chdir(p)
}

// SetOwner is an alias for Chown.
func SetOwner(p string, uid, gid int) error {
	return os.Chown(p, uid, gid)
}

// Empty removes all contents of a file or directory at the specified path. If
// the path is a directory, it removes all files and subdirectories within it
// but keeps the directory itself. If the path is a file, it truncates the file
func Empty(p string) error {
	if IsDir(p) {
		return EmptyDir(p)
	} else if IsFile(p) {
		return TruncateFile(p, 0)
	} else {
		return os.ErrNotExist
	}
}

// Link creates a hard link from src to dst. If src does not exist or dst
// already exists, it returns an error.
func Link(src, dst string) error {
	return os.Link(src, dst)
}

// Symlink creates a symbolic link from oldname to newname. If oldname does not
// exist or newname already exists, it returns an error.
func Symlink(oldname, newname string) error {
	return os.Symlink(oldname, newname)
}

// Readlink returns the destination of the named symbolic link.
func Readlink(p string) (string, error) {
	return os.Readlink(p)
}

// ForceReadlink is like Readlink but ignores any errors and returns an empty
// string in such cases.
func ForceReadlink(p string) string {
	link, _ := Readlink(p)
	return link
}

// MD5 computes the MD5 hash of a file or directory at the specified path. If the
// path is a directory, it computes the hash based on the contents of all files
// within the directory recursively. It returns the hash as a hexadecimal string.
func MD5(p string) (string, error) {
	return Hash(p, md5.New())
}

// ForceMD5 is like MD5 but ignores any errors and returns an empty string in
// such cases.
func ForceMD5(p string) string {
	sum, _ := MD5(p)
	return sum
}

// SHA1 computes the SHA-1 hash of a file or directory at the specified path. If
// the path is a directory, it computes the hash based on the contents of all files
// within the directory recursively. It returns the hash as a hexadecimal string.
func SHA1(p string) (string, error) {
	return Hash(p, sha1.New())
}

// ForceSHA1 is like SHA1 but ignores any errors and returns an empty string in
// such cases.
func ForceSHA1(p string) string {
	sum, _ := SHA1(p)
	return sum
}

// SHA256 computes the SHA-256 hash of a file or directory at the specified path.
// If the path is a directory, it computes the hash based on the contents of all
// files within the directory recursively. It returns the hash as a hexadecimal
// string.
func SHA256(path string) (string, error) {
	return Hash(path, sha256.New())
}

// ForceSHA256 is like SHA256 but ignores any errors and returns an empty string in
// such cases.
func ForceSHA256(path string) string {
	sum, _ := SHA256(path)
	return sum
}

// Checksum computes the MD5 checksum of a file or directory at the specified path.
// Alias for MD5.
func Checksum(p string) (string, error) {
	return Hash(p, md5.New())
}

// ForceChecksum is like Checksum but ignores any errors and returns an empty
// string in such cases.
func ForceChecksum(p string) string {
	sum, _ := Checksum(p)
	return sum
}

// Hash computes the hash of a file or directory at the specified path using the
// provided hash.Hash implementation. If the path is a directory, it computes the
// hash based on the contents of all files within the directory recursively.
// It returns the hash as a hexadecimal string.
func Hash(p string, h hash.Hash) (string, error) {
	if IsDir(p) {
		return hashDir(p, h)
	}
	return hashFile(p, h)
}

// ForceHash is like Hash but ignores any errors and returns an empty string in
// such cases.
func ForceHash(p string, h hash.Hash) string {
	sum, _ := Hash(p, h)
	return sum
}

// Size returns the size of a file or directory at the specified path in bytes.
// If the path is a directory, it computes the total size of all files within
// the directory recursively. It returns the size in bytes.
func Size(p string) (int64, error) {
	if IsDir(p) {
		return sizeDir(p)
	}
	return sizeFile(p)
}

// ForceSize is like Size but ignores any errors and returns zero in such cases.
func ForceSize(p string) int64 {
	size, _ := Size(p)
	return size
}

// GetModTime returns the modification time of a file at the specified path as a
// Unix timestamp (seconds since January 1, 1970). If the path does not exist
// or is a directory, it returns an error.
func GetModTime(p string) (time.Time, error) {
	info, err := os.Stat(p)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

// ForceGetModTime is like GetModTime but ignores any errors and returns the
// zero time in such cases.
func ForceGetModTime(p string) time.Time {
	t, _ := GetModTime(p)
	return t
}

// GetInfo returns a FileInfo describing the file at the specified path. If the
// path does not exist, it returns an error.
func GetInfo(p string) (os.FileInfo, error) {
	return os.Stat(p)
}

// GetMode returns the file mode (permissions) of a file at the specified path. If
// the path does not exist, it returns an error.
func GetMode(p string) (os.FileMode, error) {
	info, err := os.Stat(p)
	if err != nil {
		return 0, err
	}
	return info.Mode(), nil
}

// Getwd returns the current working directory.
func Getwd() (string, error) {
	return GetCurrentDir()
}

// Pwd is an alias for Getwd.
func Pwd() (string, error) {
	return GetCurrentDir()
}

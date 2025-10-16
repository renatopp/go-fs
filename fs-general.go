package fs

import (
	"hash"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

// ----------------------------------------------------------------------------
// CHECKS
// ----------------------------------------------------------------------------

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
		return IsEmptyDir(p)
	} else if IsFile(p) {
		return IsEmptyFile(p)
	} else {
		return true, os.ErrNotExist
	}
}

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

// IsExecutable checks if a file at the specified path is executable.
func IsExecutable(p string) bool {
	return IsExecutableFile(p)
}

// IsReadable checks if a file at the specified path is readable.
func IsReadable(p string) bool {
	return IsReadableFile(p)
}

// IsWritable checks if a file at the specified path is writable.
func IsWritable(p string) bool {
	return IsWritableFile(p)
}

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

// ----------------------------------------------------------------------------
// TRAVERSAL
// ----------------------------------------------------------------------------

// Walk traverses the directory tree rooted at the specified path, calling the
// provided function for each file or directory encountered. The function receives
// the full path of the file or directory as its argument. If the function
// returns an error, the walk is aborted and the error is returned.
func Walk(p string, fn func(string) error) error {
	return filepath.WalkDir(p, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		return fn(path)
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

// Glob returns the names of all files matching pattern or nil if there is no
// matching file. The syntax of patterns is the same as in filepath.Match.
// The pattern may describe hierarchical names such as /usr/*/bin/ed (assuming
// the Separator is '/').
func Glob(dir string, pattern string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(dir, pattern))
	if files == nil {
		files = []string{}
	}
	r := len(dir) + len(string(os.PathSeparator))
	for i, f := range files {
		files[i] = f[r:]
	}
	return files, err
}

// ----------------------------------------------------------------------------
// OPERATIONS
// ----------------------------------------------------------------------------

// Copy copies a file or directory from src to dst. If src is a directory, it
// copies the entire directory recursively. If src is a file, it copies the file.
// If dst does not exist, it will be created. If it exists, it will be merged
// (for directories) or overwritten (for files).
func Copy(src, dst string) error {
	if IsDir(src) {
		return CopyDir(src, dst)
	}
	return CopyFile(src, dst)
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

func Hide(p string) error {
	return SetHidden(p, true)
}

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

func SetOwner(p string, uid, gid int) error {
	return os.Chown(p, uid, gid)
}

// ----------------------------------------------------------------------------
// LINKS
// ----------------------------------------------------------------------------

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

// ----------------------------------------------------------------------------
// HASHING
// ----------------------------------------------------------------------------

// Checksum computes the MD5 checksum of a file or directory at the specified
// path. If the path is a directory, it computes a combined checksum of all
// files within the directory recursively. It returns the checksum as a
// hexadecimal string.
func Checksum(p string) (string, error) {
	if IsDir(p) {
		return ChecksumDir(p)
	}
	return ChecksumFile(p)
}

func Hash(p string, h hash.Hash) (string, error) {
	if IsDir(p) {
		return HashDir(p, h)
	}
	return HashFile(p, h)
}

// ----------------------------------------------------------------------------
// INFO
// ----------------------------------------------------------------------------

// Size returns the size of a file or directory at the specified path in bytes.
// If the path is a directory, it computes the total size of all files within
// the directory recursively. It returns the size in bytes.
func Size(p string) (int64, error) {
	if IsDir(p) {
		return SizeDir(p)
	}
	return SizeFile(p)
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

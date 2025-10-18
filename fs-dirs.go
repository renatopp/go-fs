package fs

import (
	"fmt"
	"hash"
	"os"
	"path/filepath"
	"slices"
)

// isEmptyDir checks if the directory at the specified path is empty. It
// returns true if the directory is empty, false if it contains files or
// subdirectories, and an error if the path does not exist or is not a
// directory.
func isEmptyDir(p string) (bool, error) {
	if !IsDir(p) {
		return false, ErrNotDir
	}
	entries, err := os.ReadDir(p)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}

// copyDir recursively copies a directory from src to dst. If dst does not
// exist, it will be created. If it exists, it will be merged with the src.
func copyDir(src, dst string) error {
	err := os.MkdirAll(dst, 0755)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
		} else {
			err = copyFile(srcPath, dstPath)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// hashDir computes a combined hash of all files in the specified
// directory using the provided hash function. It processes files recursively in a
// deterministic order to ensure consistent results. It returns the final hash
// as a hexadecimal string.
func hashDir(p string, h hash.Hash) (string, error) {
	entries, err := List(p)
	if err != nil {
		return "", err
	}

	slices.Sort(entries)

	for _, entry := range entries {
		if IsDir(entry) {
			subDirHash, err := hashDir(filepath.Join(p, entry), h)
			if err != nil {
				return "", err
			}
			h.Write([]byte(subDirHash))
		} else {
			fileHash, err := hashFile(filepath.Join(p, entry), h)
			if err != nil {
				return "", err
			}
			h.Write([]byte(fileHash))
		}
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// sizeDir computes the total size of all files within the specified directory
// and its subdirectories. It returns the total size in bytes.
func sizeDir(p string) (int64, error) {
	if !IsDir(p) {
		return 0, ErrNotDir
	}

	var totalSize int64 = 0
	entries, err := ListFilesRecursive(p)
	if err != nil {
		return 0, err
	}
	for _, entry := range entries {
		size, err := sizeFile(filepath.Join(p, entry))
		if err != nil {
			return 0, err
		}
		totalSize += size
	}
	return totalSize, nil
}

// EmptyDir removes all contents of the directory at the specified path without
// deleting the directory itself. If the directory does not exist, it returns
// an error. If the path points to a file, it returns an error.
func EmptyDir(p string) error {
	if !IsDir(p) {
		return ErrNotDir
	}
	entries, err := os.ReadDir(p)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		entryPath := filepath.Join(p, entry.Name())
		err := os.RemoveAll(entryPath)
		if err != nil {
			return err
		}
	}
	return nil
}

// IsDir checks if the given p is a directory. If the p does not exist
// or is a file, it returns false.
func IsDir(p string) bool {
	info, err := os.Stat(p)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// ListDirs returns a slice of names of all directories within the specified
// directory path. If the directory does not exist or is not accessible, it
// returns an error. This function does not include the full paths,
// only the names of the entries.
//
// This function is not recursive; it only lists entries in the specified
// directory, not in its subdirectories.
func ListDirs(p string) ([]string, error) {
	entries, err := os.ReadDir(p)
	dirs := []string{}
	if err != nil {
		return dirs, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}
	return dirs, nil
}

// ForceListDirs is like listDirs but it ignores any errors and returns only
// the slice of directory names.
func ForceListDirs(p string) []string {
	dirs, _ := ListDirs(p)
	return dirs
}

// ListDirsRecursive returns a slice of relative paths of all directories within
// the specified directory path and its subdirectories. If the directory does
// not exist or is not accessible, it returns an error. The returned paths are
// relative to the specified directory.
//
// This function is recursive; it lists directories in the specified directory
// and all its subdirectories.
func ListDirsRecursive(p string) ([]string, error) {
	if !IsDir(p) {
		return nil, ErrNotDir
	}

	results := []string{}
	err := filepath.WalkDir(p, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			relPath, err := filepath.Rel(p, path)
			if err != nil {
				return err
			}
			if relPath != "." {
				results = append(results, relPath)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// ForceListDirsRecursive is like ListDirsRecursive but it ignores any errors
// and returns only the slice of directory paths.
func ForceListDirsRecursive(p string) []string {
	dirs, _ := ListDirsRecursive(p)
	return dirs
}

// CreateDir creates a directory at the specified path, including any necessary
// parent directories. If the directory already exists, it does nothing and
// returns nil.
func CreateDir(p string) error {
	return os.MkdirAll(p, 0755)
}

// EnsureDir ensures that a directory exists at the specified path. It follows
// these rules:
//
//   - If the p points to an existing file, it returns an error.
//   - If the directory already exists, it does nothing and returns nil.
//   - If the directory does not exist, it creates the directory along with any
//     necessary parent directories.
func EnsureDir(p string) error {
	if IsFile(p) {
		return ErrIsFile
	}
	if Exists(p) {
		return nil
	}
	return os.MkdirAll(p, 0755)
}

// CreateTempDir creates a temporary directory with the specified prefix in the
// system's default temporary directory. It returns the full path of the created
// directory.
func CreateTempDir(prefix string) (string, error) {
	return os.MkdirTemp("", prefix)
}

// ForceCreateTempDir is like CreateTempDir but it ignores any errors and
// returns only the path of the created temporary directory.
func ForceCreateTempDir(prefix string) string {
	dir, _ := CreateTempDir(prefix)
	return dir
}

// GetCurrentDir is an alias for Getwd.
func GetCurrentDir() (string, error) {
	return os.Getwd()
}

// ForceGetCurrentDir is like GetCurrentDir but it ignores any errors and
// returns only the current directory path.
func ForceGetCurrentDir() string {
	dir, _ := GetCurrentDir()
	return dir
}

// GetTempDir returns the default temporary directory of the system.
func GetTempDir() string {
	return os.TempDir()
}

// GetCacheDir returns the cache directory of the current user.
func GetCacheDir() (string, error) {
	return os.UserCacheDir()
}

// ForceGetCacheDir is like GetCacheDir but it ignores any errors and
// returns only the cache directory path.
func ForceGetCacheDir() string {
	dir, _ := GetCacheDir()
	return dir
}

// GetConfigDir returns the configuration directory of the current user.
func GetConfigDir() (string, error) {
	return os.UserConfigDir()
}

// ForceGetConfigDir is like GetConfigDir but it ignores any errors and
// returns only the configuration directory path.
func ForceGetConfigDir() string {
	dir, _ := GetConfigDir()
	return dir
}

// GetHomeDir returns the home directory of the current user.
func GetHomeDir() (string, error) {
	return os.UserHomeDir()
}

// ForceGetHomeDir is like GetHomeDir but it ignores any errors and
// returns only the home directory path.
func ForceGetHomeDir() string {
	dir, _ := GetHomeDir()
	return dir
}

// GetParentDir returns the parent directory of the given path p. If p is a file,
// it returns the parent directory of that file. If p is a directory, it returns
// that directory itself. If the path does not exist, it returns an error but
// also returns the parent directory path treating p as a file.
func GetParentDir(p string) (string, error) {
	p = Force(AbsolutePath(p))
	if IsFile(p) {
		return GetPathParent(p), nil
	}
	if IsDir(p) {
		return p, nil
	}
	return GetPathParent(p), ErrNotExist
}

// ForceGetParentDir is like GetParentDir but it ignores any errors and
// returns only the parent directory path.
func ForceGetParentDir(p string) string {
	dir, _ := GetParentDir(p)
	return dir
}

// GetParentDirName returns the name of the parent directory of the given path p.
// If p is a file, it returns the name of its parent directory. If p is a
// directory, it returns the name of that directory itself. If the path does not
// exist, it returns an error but also returns the parent directory name treating p
// as a file.
func GetParentDirName(p string) (string, error) {
	p = Force(AbsolutePath(p))
	if IsFile(p) {
		return GetPathParentName(p), nil
	}
	if IsDir(p) {
		return GetPathBase(p), nil
	}
	return GetPathParentName(p), ErrNotExist
}

// ForceGetParentDirName is like GetParentDirName but it ignores any errors and
// returns only the parent directory name.
func ForceGetParentDirName(p string) string {
	name, _ := GetParentDirName(p)
	return name
}

// GetDirParts returns the PathParts structure for the given directory path p.
// If p is a file, it returns the parts of its parent directory.
func GetDirParts(p string) PathParts {
	p = Force(AbsolutePath(p))
	if IsDir(p) {
		return PathParts{
			Absolute:   p,
			Base:       GetPathBase(p),
			Name:       GetPathBase(p),
			Ext:        "",
			ExtName:    "",
			Parent:     GetPathParent(p),
			ParentName: GetPathParentName(p),
			Volume:     GetPathVolume(p),
		}
	}
	return GetPathParts(p)
}

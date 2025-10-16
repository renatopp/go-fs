package fs

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"hash"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// ----------------------------------------------------------------------------
// CHECKS
// ----------------------------------------------------------------------------

// IsFile checks if the given p is a file. If the p does not exist or is
// a directory, it returns false.
func IsFile(p string) bool {
	info, err := os.Stat(p)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// IsSameFile checks if two paths point to the same file.
func IsSameFile(p1, p2 string) bool {
	if !IsFile(p1) || !IsFile(p2) {
		return false
	}
	return IsSame(p1, p2)
}

// IsEmptyFile checks if the file at the specified path is empty (has a size of
// zero bytes). If the path points to a directory or does not exist, it returns
// an error.
func IsEmptyFile(p string) (bool, error) {
	info, err := os.Stat(p)
	if err != nil {
		return false, err
	}
	if info.IsDir() {
		return false, ErrIsDir
	}
	return info.Size() == 0, nil
}

// IsExecutable checks if a file at the specified path is executable.
func IsExecutableFile(p string) bool {
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

// IsReadable checks if a file at the specified path is readable.
func IsReadableFile(p string) bool {
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

// IsWritable checks if a file at the specified path is writable.
func IsWritableFile(p string) bool {
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

func IsHiddenFile(p string) (bool, error) {
	if !IsFile(p) {
		return false, ErrNotFile
	}
	return IsHidden(p)
}

// ----------------------------------------------------------------------------
// TRAVERSAL
// ----------------------------------------------------------------------------

// ListFiles returns a slice of names of all files within the specified
// directory path. If the directory does not exist or is not accessible, it
// returns an error. This function does not include the full paths,
// only the names of the entries.
//
// This function is not recursive; it only lists entries in the specified
// directory, not in its subdirectories.
func ListFiles(p string) ([]string, error) {
	entries, err := os.ReadDir(p)
	files := []string{}
	if err != nil {
		return files, err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

// ListFilesRecursive returns a slice of relative paths of all files within the
// specified directory path and its subdirectories. If the directory does not
// exist or is not accessible, it returns an error. The returned paths are
// relative to the specified directory.
//
// This function is recursive; it lists files in the specified directory
// and all its subdirectories.
func ListFilesRecursive(p string) ([]string, error) {
	if !IsDir(p) {
		return nil, ErrNotDir
	}

	results := []string{}
	err := filepath.WalkDir(p, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			relPath, err := filepath.Rel(p, path)
			if err != nil {
				return err
			}
			results = append(results, relPath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// ----------------------------------------------------------------------------
// OPERATIONS
// ----------------------------------------------------------------------------

// ReadFile reads the entire content of a file and returns it as a byte slice.
func ReadFile(p string) ([]byte, error) {
	return os.ReadFile(p)
}

// ReadFileString reads the entire content of a file and returns it as a string.
func ReadFileString(p string) (string, error) {
	data, err := ReadFile(p)
	return string(data), err
}

// ReadFileLines reads a file and returns its content as a slice of strings,
// where each string represents a line in the file.
func ReadFileLines(p string) ([]string, error) {
	data, err := ReadFile(p)
	if err != nil {
		return []string{}, err
	}
	return strings.Split(string(data), "\n"), nil
}

// ReadFileJson reads a JSON file and unmarshals its content into the provided
// variable v, which should be a pointer to the desired data structure.
func ReadFileJson(p string, v any) error {
	data, err := ReadFile(p)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// WriteFile writes the given byte slice data to a file at the specified path.
// IF the directory does not exist, it will fail. If the file exists, it will
// be overwritten.
func WriteFile(p string, data []byte) error {
	return os.WriteFile(p, data, 0644)
}

// WriteFileString writes the given string data to a file at the specified path.
// If the directory does not exist, it will fail. If the file exists, it will
// be overwritten.
func WriteFileString(p string, data string) error {
	return WriteFile(p, []byte(data))
}

// WriteFileLines writes the given slice of strings to a file at the specified
// path, with each string representing a line in the file. If the directory
// does not exist, it will fail. If the file exists, it will be overwritten.
func WriteFileLines(p string, lines []string) error {
	data := strings.Join(lines, "\n")
	return WriteFileString(p, data)
}

// WriteFileJson marshals the given variable v into JSON format and writes it to
// a file at the specified path. If the directory does not exist, it will fail.
// If the file exists, it will be overwritten.
func WriteFileJson(p string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return WriteFile(p, data)
}

// AppendFile appends the given byte slice data to a file at the specified path.
// If the file does not exist, it will be created.
func AppendFile(p string, data []byte) error {
	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

// AppendFileString appends the given string data to a file at the specified path.
// If the file does not exist, it will be created.
func AppendFileString(p string, data string) error {
	return AppendFile(p, []byte(data))
}

// AppendFileLines appends the given slice of strings to a file at the specified
// p, with each string representing a line in the file. If the file does not
// exist, it will be created.
// If the file exists, a newline will be added before appending the new lines.
func AppendFileLines(p string, lines []string) error {
	data := strings.Join(lines, "\n")
	if Exists(p) {
		data = "\n" + data
	}
	return AppendFileString(p, data)
}

// AppendFileJson appends the JSON representation of the given variable v to a
// file at the specified path. If the file does not exist, it will be created.
// Json will be appended without indentation or newlines.
// If the file exists, a newline will be added before appending the new JSON.
func AppendFileJson(p string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	str := string(data)
	if Exists(p) {
		str = "\n" + str
	}
	return AppendFile(p, []byte(str))
}

// TouchFile creates an empty file at the specified path if it does not already
// exist.
func TouchFile(p string) error {
	if Exists(p) {
		return nil
	}
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	return f.Close()
}

// EnsureFile ensures that a file exists at the specified path. It follows
// these rules:
//
//   - If the p points to an existing directory, it returns an error.
//   - If the file already exists, it does nothing and returns nil.
//   - If the file does not exist, it creates any necessary parent directories
//     and then creates an empty file at the specified path.
func EnsureFile(p string) error {
	if IsDir(p) {
		return ErrIsDir
	}
	if Exists(p) {
		return nil
	}

	dir := filepath.Dir(p)
	err := EnsureDir(dir)
	if err != nil {
		return err
	}

	return TouchFile(p)
}

// ReplaceInFile reads the content of the file at the specified path, replaces
// all occurrences of the old byte slice with the new byte slice, and writes
// the modified content back to the file. If the old byte slice is not found
// in the file, it does nothing.
func ReplaceInFile(p string, old []byte, new []byte) error {
	data, err := ReadFile(p)
	if err != nil {
		return err
	}
	if !bytes.Contains(data, old) {
		return nil
	}
	modified := bytes.ReplaceAll(data, old, new)
	return WriteFile(p, modified)
}

// ReplaceInFileString is like ReplaceInFile but works with strings instead of
// byte slices.
func ReplaceInFileString(p string, old string, new string) error {
	return ReplaceInFile(p, []byte(old), []byte(new))
}

// CreateTempFile creates a temporary file with the specified prefix in the system's
// default temporary directory. It returns the full path of the created file.
func CreateTempFile(prefix string) (string, error) {
	f, err := os.CreateTemp("", prefix)
	if err != nil {
		return "", err
	}
	f.Close()
	return f.Name(), nil
}

// CreateTempFileOpen creates a temporary file with the specified prefix in the
// system's default temporary directory and returns an open file handle to it.
func CreateTempFileOpen(prefix string) (*os.File, error) {
	return os.CreateTemp("", prefix)
}

// CopyFile copies a file from src to dst. If dst does not exist, it will be
// created. If it exists, it will be overwritten.
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = srcFile.WriteTo(dstFile)
	return err
}

// MoveFile moves a file from src to dst. If dst exists, it will be overwritten.
func MoveFile(src, dst string) error {
	if !IsFile(src) {
		return ErrNotFile
	}
	return os.Rename(src, dst)
}

// RenameFile renames a file from oldPath to newPath. It is an alias for MoveFile.
func RenameFile(oldPath, newPath string) error {
	return MoveFile(oldPath, newPath)
}

// DeleteFile deletes the file at the specified path. If the path points to a
// directory or does not exist, it returns an error.
func DeleteFile(p string) error {
	if !IsFile(p) {
		return ErrNotFile
	}
	return os.Remove(p)
}

// TruncateFile truncates the file at the specified path to the given size in
// bytes. If the path points to a directory or does not exist, it returns an
// error.
func TruncateFile(p string, size int64) error {
	if !IsFile(p) {
		return ErrIsDir
	}
	return os.Truncate(p, size)
}

func SetModeFile(p string, mode os.FileMode) error {
	if !IsFile(p) {
		return ErrNotFile
	}
	return os.Chmod(p, mode)
}

func SetHiddenFile(p string, hidden bool) error {
	if !IsFile(p) {
		return ErrNotFile
	}
	return SetHidden(p, hidden)
}

func HideFile(p string) error {
	return SetHiddenFile(p, true)
}

func UnhideFile(p string) error {
	return SetHiddenFile(p, false)
}

func SetOwnerFile(p string, uid, gid int) error {
	if !IsFile(p) {
		return ErrNotFile
	}
	return SetOwner(p, uid, gid)
}

// ----------------------------------------------------------------------------
// LINKS
// ----------------------------------------------------------------------------

// LinkFile creates a hard link from oldname to newname. If oldname does not
// exist or is not a file, it returns an error.
func LinkFile(oldname, newname string) error {
	if !IsFile(oldname) {
		return ErrNotFile
	}
	return Link(oldname, newname)
}

// SymlinkFile creates a symbolic link from oldname to newname. If oldname
// does not exist or is not a file, it returns an error.
func SymlinkFile(oldname, newname string) error {
	if !IsFile(oldname) {
		return ErrNotFile
	}
	return Symlink(oldname, newname)
}

// ReadlinkFile reads the target of a symbolic link. If the path does not
// point to a file or is not a symbolic link, it returns an error.
func ReadlinkFile(p string) (string, error) {
	if !IsFile(p) {
		return "", ErrNotFile
	}
	return Readlink(p)
}

// ----------------------------------------------------------------------------
// HASHING
// ----------------------------------------------------------------------------

// MD5File computes the MD5 checksum of the file at the specified path and
// returns it as a hexadecimal string.
func MD5File(p string) (string, error) {
	return HashFile(p, md5.New())
}

// SHA1File computes the SHA-1 checksum of the file at the specified path and
// returns it as a hexadecimal string.
func SHA1File(p string) (string, error) {
	return HashFile(p, sha1.New())
}

// SHA256File computes the SHA-256 checksum of the file at the specified path and
// returns it as a hexadecimal string.
func SHA256File(p string) (string, error) {
	return HashFile(p, sha256.New())
}

// ChecksumFile computes the MD5 checksum of the file at the specified path and
// returns it as a hexadecimal string. It is an alias for MD5File.
func ChecksumFile(p string) (string, error) {
	return MD5File(p)
}

// HashFile computes the hash of the file at the specified path using the
// provided hash function and returns it as a hexadecimal string.
func HashFile(p string, h hash.Hash) (string, error) {
	data, err := ReadFile(p)
	if err != nil {
		return "", err
	}
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// ----------------------------------------------------------------------------
// INFO
// ----------------------------------------------------------------------------

// SizeFile returns the size of the file at the specified path in bytes. If the
// path points to a directory or does not exist, it returns an error.
func SizeFile(p string) (int64, error) {
	info, err := os.Stat(p)
	if err != nil {
		return 0, err
	}
	if info.IsDir() {
		return 0, ErrIsDir
	}
	return info.Size(), nil
}

func GetModTimeFile(p string) (time.Time, error) {
	if !IsFile(p) {
		return time.Time{}, ErrNotFile
	}
	return GetModTime(p)
}

func GetInfoFile(p string) (os.FileInfo, error) {
	if !IsFile(p) {
		return nil, ErrNotFile
	}
	return GetInfo(p)
}

func GetModeFile(p string) (os.FileMode, error) {
	if !IsFile(p) {
		return 0, ErrNotFile
	}
	return GetMode(p)
}

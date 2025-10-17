package fs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash"
	"os"
	"path/filepath"
	"strings"
)

// ----------------------------------------------------------------------------
// INTERNAL
// ----------------------------------------------------------------------------

// isEmptyFile checks if the file at the specified path is empty.
//
// It returns a value and an error. The value is true if the file is empty,
// false if it contains data or if there was an error. If the path points to a
// directory or does not exist, it returns an appropriate error.
func isEmptyFile(p string) (bool, error) {
	info, err := os.Stat(p)
	if err != nil {
		return false, err
	}
	if info.IsDir() {
		return false, ErrIsDir
	}
	return info.Size() == 0, nil
}

// copyFile copies a file from src to dst. If dst does not exist, it will be
// created. If it exists, it will be overwritten.
func copyFile(src, dst string) error {
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

// sizeFile returns the size of the file at the specified path in bytes. If the
// path points to a directory or does not exist, it returns an error.
func sizeFile(p string) (int64, error) {
	info, err := os.Stat(p)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// hashFile computes the hash of the file at the specified path using the
// provided hash function and returns it as a hexadecimal string.
func hashFile(p string, h hash.Hash) (string, error) {
	data, err := ReadFile(p)
	if err != nil {
		return "", err
	}
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// ----------------------------------------------------------------------------
// PUBLIC
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

func ForceListFiles(p string) []string {
	files, _ := ListFiles(p)
	return files
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

func ForceListFilesRecursive(p string) []string {
	files, _ := ListFilesRecursive(p)
	return files
}

// ReadFile reads the entire content of a file and returns it as a byte slice.
func ReadFile(p string) ([]byte, error) {
	return os.ReadFile(p)
}

func ForceReadFile(p string) []byte {
	data, err := ReadFile(p)
	if err != nil {
		return []byte{}
	}
	return data
}

// ReadFileString reads the entire content of a file and returns it as a string.
func ReadFileString(p string) (string, error) {
	data, err := ReadFile(p)
	return string(data), err
}

func ForceReadFileString(p string) string {
	str, _ := ReadFileString(p)
	return str
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

func ForceReadFileLines(p string) []string {
	lines, _ := ReadFileLines(p)
	return lines
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

func ForceCreateTempFile(prefix string) string {
	p, _ := CreateTempFile(prefix)
	return p
}

// CreateTempFileOpen creates a temporary file with the specified prefix in the
// system's default temporary directory and returns an open file handle to it.
func CreateTempFileOpen(prefix string) (*os.File, error) {
	return os.CreateTemp("", prefix)
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

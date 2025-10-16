package fs

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"strings"
)

var ErrIsDir = errors.New("is a directory")

// Exists checks if a file or directory exists at the given path.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsFile checks if the given path is a file. If the path does not exist or is
// a directory, it returns false.
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// IsDir checks if the given path is a directory. If the path does not exist
// or is a file, it returns false.
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// ----------------------------------------------------------------------------
// FILE-RELATED OPERATIONS
// ----------------------------------------------------------------------------

// ReadFile reads the entire content of a file and returns it as a byte slice.
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// ReadFileString reads the entire content of a file and returns it as a string.
func ReadFileString(path string) (string, error) {
	data, err := ReadFile(path)
	return string(data), err
}

// ReadFileLines reads a file and returns its content as a slice of strings,
// where each string represents a line in the file.
func ReadFileLines(path string) ([]string, error) {
	data, err := ReadFile(path)
	if err != nil {
		return []string{}, err
	}
	return strings.Split(string(data), "\n"), nil
}

// ReadFileJson reads a JSON file and unmarshals its content into the provided
// variable v, which should be a pointer to the desired data structure.
func ReadFileJson(path string, v any) error {
	data, err := ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// WriteFile writes the given byte slice data to a file at the specified path.
// IF the directory does not exist, it will fail. If the file exists, it will
// be overwritten.
func WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

// WriteFileString writes the given string data to a file at the specified path.
// If the directory does not exist, it will fail. If the file exists, it will
// be overwritten.
func WriteFileString(path string, data string) error {
	return WriteFile(path, []byte(data))
}

// WriteFileLines writes the given slice of strings to a file at the specified
// path, with each string representing a line in the file. If the directory
// does not exist, it will fail. If the file exists, it will be overwritten.
func WriteFileLines(path string, lines []string) error {
	data := strings.Join(lines, "\n")
	return WriteFileString(path, data)
}

// WriteFileJson marshals the given variable v into JSON format and writes it to
// a file at the specified path. If the directory does not exist, it will fail.
// If the file exists, it will be overwritten.
func WriteFileJson(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return WriteFile(path, data)
}

// AppendFile appends the given byte slice data to a file at the specified path.
// If the file does not exist, it will be created.
func AppendFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

// AppendFileString appends the given string data to a file at the specified path.
// If the file does not exist, it will be created.
func AppendFileString(path string, data string) error {
	return AppendFile(path, []byte(data))
}

// AppendFileLines appends the given slice of strings to a file at the specified
// path, with each string representing a line in the file. If the file does not
// exist, it will be created.
// If the file exists, a newline will be added before appending the new lines.
func AppendFileLines(path string, lines []string) error {
	data := strings.Join(lines, "\n")
	if Exists(path) {
		data = "\n" + data
	}
	return AppendFileString(path, data)
}

// AppendFileJson appends the JSON representation of the given variable v to a
// file at the specified path. If the file does not exist, it will be created.
// Json will be appended without indentation or newlines.
// If the file exists, a newline will be added before appending the new JSON.
func AppendFileJson(path string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	str := string(data)
	if Exists(path) {
		str = "\n" + str
	}
	return AppendFile(path, []byte(str))
}

// EnsureFile ensures that a file exists at the specified path. If the file does
// not exist, it creates an empty file. If the directories in the path do not
// exist, they will be created as well. If the file already exists, it does nothing,
// but if the path is a directory, it will return an error.
func EnsureFile(p string) error {
	if IsDir(p) {
		return ErrIsDir
	}

	dirs := path.Dir(p)
	err := os.MkdirAll(dirs, os.ModePerm)
	if err != nil {
		return err
	}

	if Exists(p) {
		return nil
	}
	file, err := os.Create(p)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

// // File operations for raw
// func EnsureFile(path string) error
// func ReplaceInFile(path string, old []byte, new []byte) error
// func Touch(path string) error
// func TempFile(prefix string) (string, error)
// func IsSameFile(path1, path2 string) bool

// // Directory operations
// func CreateDir(path string) error
// func EmptyDir(path string) error
// func EnsureDir(path string) error
// func List(path string) ([]string, error)
// func ListFiles(path string) ([]string, error)
// func ListDirs(path string) ([]string, error)
// func ListRecursive(path string) ([]string, error)
// func Glob(pattern string) ([]string, error)
// func TempDir(prefix string) (string, error)
// func Watch(path string, recursive bool, callback func(event string, path string)) error
// func WatchGlob(pattern string, callback func(event string, path string)) error

// // Agnostic
// func Copy(src, dst string) error
// func Move(src, dst string) error
// func Remove(path string) error
// func Checksum(path string) (string, error)
// func Size(path string) int64

// // Path operations
// func Join(elem ...string) string
// func JoinLinux(elem ...string) string
// func JoinWindows(elem ...string) string
// func JoinWith(sep string, elem ...string) string
// func Abs(path string) (string, error)
// func Base(path string) string
// func Dir(path string) string
// func Ext(path string) string
// func HasExt(path string) bool
// func Stem(path string) string
// //...

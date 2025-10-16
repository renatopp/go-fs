package fs_test

import (
	"os"
	"testing"

	"github.com/renatopp/go-fs"
	"github.com/stretchr/testify/assert"
)

func failIf(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func tempFile(t *testing.T) (*os.File, func()) {
	tmpFile, err := os.CreateTemp("", "test")
	failIf(t, err)
	return tmpFile, func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}
}

func tempFileWithContent(t *testing.T, content string) (*os.File, func()) {
	tmpFile, err := os.CreateTemp("", "test")
	failIf(t, err)
	_, err = tmpFile.WriteString(content)
	failIf(t, err)
	return tmpFile, func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}
}

func tempDir(t *testing.T) (string, func()) {
	tmpDir, err := os.MkdirTemp("", "test")
	failIf(t, err)
	return tmpDir, func() {
		os.RemoveAll(tmpDir)
	}
}

func Test_Exists(t *testing.T) {
	file, fileDone := tempFile(t)
	defer fileDone()
	assert.True(t, fs.Exists(file.Name()))
	assert.False(t, fs.Exists(file.Name()+"-nonexistent"))

	dir, dirDone := tempDir(t)
	defer dirDone()
	assert.True(t, fs.Exists(dir))
	assert.False(t, fs.Exists(dir+"-nonexistent"))
}

func Test_IsFile(t *testing.T) {
	file, fileDone := tempFile(t)
	defer fileDone()
	assert.True(t, fs.IsFile(file.Name()))
	assert.False(t, fs.IsFile(file.Name()+"-nonexistent"))

	dir, dirDone := tempDir(t)
	defer dirDone()
	assert.False(t, fs.IsFile(dir))
	assert.False(t, fs.IsFile(dir+"-nonexistent"))
}

func Test_IsDir(t *testing.T) {
	file, fileDone := tempFile(t)
	defer fileDone()
	assert.False(t, fs.IsDir(file.Name()))
	assert.False(t, fs.IsDir(file.Name()+"-nonexistent"))

	dir, dirDone := tempDir(t)
	defer dirDone()
	assert.True(t, fs.IsDir(dir))
	assert.False(t, fs.IsDir(dir+"-nonexistent"))
}

//

func Test_ReadFile(t *testing.T) {
	content := "Hello, World!"
	file, fileDone := tempFileWithContent(t, content)
	defer fileDone()

	data, err := fs.ReadFile(file.Name())
	failIf(t, err)
	assert.Equal(t, []byte(content), data)
}

func Test_ReadFileString(t *testing.T) {
	content := "Hello, World!"
	file, fileDone := tempFileWithContent(t, content)
	defer fileDone()

	data, err := fs.ReadFileString(file.Name())
	failIf(t, err)
	assert.Equal(t, content, data)
}

func Test_ReadFileLines(t *testing.T) {
	file, fileDone := tempFileWithContent(t, "Line 1\nLine 2\nLine 3")
	defer fileDone()

	lines, err := fs.ReadFileLines(file.Name())
	failIf(t, err)
	assert.Equal(t, []string{"Line 1", "Line 2", "Line 3"}, lines)

	lines2, err2 := fs.ReadFileLines(file.Name() + "-nonexistent")
	assert.Error(t, err2)
	assert.Equal(t, []string{}, lines2)
}

func Test_ReadFileJson(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	content := `{"name": "John", "age": 30}`
	file, fileDone := tempFileWithContent(t, content)
	defer fileDone()

	var result Person
	err := fs.ReadFileJson(file.Name(), &result)
	failIf(t, err)
	assert.Equal(t, "John", result.Name)
	assert.Equal(t, 30, result.Age)

	var result2 Person
	err2 := fs.ReadFileJson(file.Name()+"-nonexistent", &result2)
	assert.Error(t, err2)
	assert.Equal(t, Person{}, result2)
}

func Test_WriteFile(t *testing.T) {
	content := "Hello, WriteFile!"
	file, fileDone := tempFile(t)
	defer fileDone()

	err := fs.WriteFile(file.Name(), []byte(content))
	failIf(t, err)

	data, err := fs.ReadFileString(file.Name())
	failIf(t, err)
	assert.Equal(t, content, data)
}

func Test_WriteFileString(t *testing.T) {
	content := "Hello, WriteFileString!"
	file, fileDone := tempFile(t)
	defer fileDone()

	err := fs.WriteFileString(file.Name(), content)
	failIf(t, err)

	data, err := fs.ReadFileString(file.Name())
	failIf(t, err)
	assert.Equal(t, content, data)
}

func Test_WriteFileLines(t *testing.T) {
	lines := []string{"Line 1", "Line 2", "Line 3"}
	file, fileDone := tempFile(t)
	defer fileDone()

	err := fs.WriteFileLines(file.Name(), lines)
	failIf(t, err)

	data, err := fs.ReadFileLines(file.Name())
	failIf(t, err)
	assert.Equal(t, lines, data)
}

func Test_WriteFileJson(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	person := Person{Name: "Alice", Age: 25}
	file, fileDone := tempFile(t)
	defer fileDone()

	err := fs.WriteFileJson(file.Name(), person)
	failIf(t, err)

	var result Person
	err = fs.ReadFileJson(file.Name(), &result)
	failIf(t, err)
	assert.Equal(t, person, result)
}

func Test_WriteFileJson_Invalid(t *testing.T) {
	ch := make(chan int) // channels cannot be JSON marshaled
	file, fileDone := tempFile(t)
	defer fileDone()

	err := fs.WriteFileJson(file.Name(), ch)
	assert.Error(t, err)
}

func Test_AppendFile(t *testing.T) {
	content1 := "Hello"
	content2 := ", World!"
	file, fileDone := tempFile(t)
	defer fileDone()

	err := fs.AppendFile(file.Name(), []byte(content1))
	failIf(t, err)

	err = fs.AppendFile(file.Name(), []byte(content2))
	failIf(t, err)

	data, err := fs.ReadFileString(file.Name())
	failIf(t, err)
	assert.Equal(t, content1+content2, data)
}

func Test_AppendFile_Invalid(t *testing.T) {
	err := fs.AppendFile("", nil)
	assert.Error(t, err)
}

func Test_AppendFileString(t *testing.T) {
	content1 := "Hello"
	content2 := ", World!"
	file, fileDone := tempFile(t)
	defer fileDone()
	err := fs.AppendFileString(file.Name(), content1)
	failIf(t, err)

	err = fs.AppendFileString(file.Name(), content2)
	failIf(t, err)

	data, err := fs.ReadFileString(file.Name())
	failIf(t, err)
	assert.Equal(t, content1+content2, data)
}

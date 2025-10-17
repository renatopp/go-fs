package fs_test

// func failIf(t *testing.T, err error) {
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// func tempFile(t *testing.T) (*os.File, func()) {
// 	tmpFile, err := os.CreateTemp("", "test")
// 	failIf(t, err)
// 	return tmpFile, func() {
// 		tmpFile.Close()
// 		os.Remove(tmpFile.Name())
// 	}
// }

// func tempFileWithContent(t *testing.T, content string) (*os.File, func()) {
// 	tmpFile, err := os.CreateTemp("", "test")
// 	failIf(t, err)
// 	_, err = tmpFile.WriteString(content)
// 	failIf(t, err)
// 	return tmpFile, func() {
// 		tmpFile.Close()
// 		os.Remove(tmpFile.Name())
// 	}
// }

// func tempDir(t *testing.T) (string, func()) {
// 	tmpDir, err := os.MkdirTemp("", "test")
// 	failIf(t, err)
// 	return tmpDir, func() {
// 		os.RemoveAll(tmpDir)
// 	}
// }

// // FILES

// func Test_IsFile(t *testing.T) {
// 	file, fileDone := tempFile(t)
// 	defer fileDone()
// 	assert.True(t, fs.IsFile(file.Name()))
// 	assert.False(t, fs.IsFile(file.Name()+"-nonexistent"))

// 	dir, dirDone := tempDir(t)
// 	defer dirDone()
// 	assert.False(t, fs.IsFile(dir))
// 	assert.False(t, fs.IsFile(dir+"-nonexistent"))
// }

// func Test_ReadFile(t *testing.T) {
// 	content := "Hello, World!"
// 	file, fileDone := tempFileWithContent(t, content)
// 	defer fileDone()

// 	data, err := fs.ReadFile(file.Name())
// 	failIf(t, err)
// 	assert.Equal(t, []byte(content), data)
// }

// func Test_ReadFileString(t *testing.T) {
// 	content := "Hello, World!"
// 	file, fileDone := tempFileWithContent(t, content)
// 	defer fileDone()

// 	data, err := fs.ReadFileString(file.Name())
// 	failIf(t, err)
// 	assert.Equal(t, content, data)
// }

// func Test_ReadFileLines(t *testing.T) {
// 	file, fileDone := tempFileWithContent(t, "Line 1\nLine 2\nLine 3")
// 	defer fileDone()

// 	lines, err := fs.ReadFileLines(file.Name())
// 	failIf(t, err)
// 	assert.Equal(t, []string{"Line 1", "Line 2", "Line 3"}, lines)

// 	lines2, err2 := fs.ReadFileLines(file.Name() + "-nonexistent")
// 	assert.Error(t, err2)
// 	assert.Equal(t, []string{}, lines2)
// }

// func Test_ReadFileJson(t *testing.T) {
// 	type Person struct {
// 		Name string `json:"name"`
// 		Age  int    `json:"age"`
// 	}
// 	content := `{"name": "John", "age": 30}`
// 	file, fileDone := tempFileWithContent(t, content)
// 	defer fileDone()

// 	var result Person
// 	err := fs.ReadFileJson(file.Name(), &result)
// 	failIf(t, err)
// 	assert.Equal(t, "John", result.Name)
// 	assert.Equal(t, 30, result.Age)

// 	var result2 Person
// 	err2 := fs.ReadFileJson(file.Name()+"-nonexistent", &result2)
// 	assert.Error(t, err2)
// 	assert.Equal(t, Person{}, result2)
// }

// func Test_WriteFile(t *testing.T) {
// 	content := "Hello, WriteFile!"
// 	file, fileDone := tempFile(t)
// 	defer fileDone()

// 	err := fs.WriteFile(file.Name(), []byte(content))
// 	failIf(t, err)

// 	data, err := fs.ReadFileString(file.Name())
// 	failIf(t, err)
// 	assert.Equal(t, content, data)
// }

// func Test_WriteFileString(t *testing.T) {
// 	content := "Hello, WriteFileString!"
// 	file, fileDone := tempFile(t)
// 	defer fileDone()

// 	err := fs.WriteFileString(file.Name(), content)
// 	failIf(t, err)

// 	data, err := fs.ReadFileString(file.Name())
// 	failIf(t, err)
// 	assert.Equal(t, content, data)
// }

// func Test_WriteFileLines(t *testing.T) {
// 	lines := []string{"Line 1", "Line 2", "Line 3"}
// 	file, fileDone := tempFile(t)
// 	defer fileDone()

// 	err := fs.WriteFileLines(file.Name(), lines)
// 	failIf(t, err)

// 	data, err := fs.ReadFileLines(file.Name())
// 	failIf(t, err)
// 	assert.Equal(t, lines, data)
// }

// func Test_WriteFileJson(t *testing.T) {
// 	type Person struct {
// 		Name string `json:"name"`
// 		Age  int    `json:"age"`
// 	}
// 	person := Person{Name: "Alice", Age: 25}
// 	file, fileDone := tempFile(t)
// 	defer fileDone()

// 	err := fs.WriteFileJson(file.Name(), person)
// 	failIf(t, err)

// 	var result Person
// 	err = fs.ReadFileJson(file.Name(), &result)
// 	failIf(t, err)
// 	assert.Equal(t, person, result)
// }

// func Test_WriteFileJson_Invalid(t *testing.T) {
// 	ch := make(chan int) // channels cannot be JSON marshaled
// 	file, fileDone := tempFile(t)
// 	defer fileDone()

// 	err := fs.WriteFileJson(file.Name(), ch)
// 	assert.Error(t, err)
// }

// func Test_AppendFile(t *testing.T) {
// 	content1 := "Hello"
// 	content2 := ", World!"
// 	file, fileDone := tempFile(t)
// 	defer fileDone()

// 	err := fs.AppendFile(file.Name(), []byte(content1))
// 	failIf(t, err)

// 	err = fs.AppendFile(file.Name(), []byte(content2))
// 	failIf(t, err)

// 	data, err := fs.ReadFileString(file.Name())
// 	failIf(t, err)
// 	assert.Equal(t, content1+content2, data)
// }

// func Test_AppendFile_Invalid(t *testing.T) {
// 	err := fs.AppendFile("", nil)
// 	assert.Error(t, err)
// }

// func Test_AppendFileString(t *testing.T) {
// 	content1 := "Hello"
// 	content2 := ", World!"
// 	file, fileDone := tempFile(t)
// 	defer fileDone()
// 	err := fs.AppendFileString(file.Name(), content1)
// 	failIf(t, err)

// 	err = fs.AppendFileString(file.Name(), content2)
// 	failIf(t, err)

// 	data, err := fs.ReadFileString(file.Name())
// 	failIf(t, err)
// 	assert.Equal(t, content1+content2, data)
// }

// func Test_AppendFileLines(t *testing.T) {
// 	lines1 := []string{"Line 1", "Line 2"}
// 	lines2 := []string{"Line 3", "Line 4"}
// 	tempName := filepath.Join(os.TempDir(), "appendlines")
// 	defer os.RemoveAll(tempName)

// 	err := fs.AppendFileLines(tempName, lines1)
// 	failIf(t, err)

// 	err = fs.AppendFileLines(tempName, lines2)
// 	failIf(t, err)

// 	data, err := fs.ReadFileLines(tempName)
// 	failIf(t, err)
// 	assert.Equal(t, append(lines1, lines2...), data)
// }

// func Test_AppendFileJson(t *testing.T) {
// 	type Person struct {
// 		Name string `json:"name"`
// 	}
// 	person1 := Person{Name: "Alice"}
// 	person2 := Person{Name: "Bob"}
// 	tempName := filepath.Join(os.TempDir(), "appendlines")
// 	defer os.RemoveAll(tempName)

// 	err := fs.AppendFileJson(tempName, person1)
// 	failIf(t, err)

// 	err = fs.AppendFileJson(tempName, person2)
// 	failIf(t, err)

// 	lines, err := fs.ReadFileLines(tempName)
// 	failIf(t, err)
// 	assert.Equal(t, "{\"name\":\"Alice\"}", lines[0])
// 	assert.Equal(t, "{\"name\":\"Bob\"}", lines[1])
// }

// func Test_AppendFileJson_Invalid(t *testing.T) {
// 	ch := make(chan int) // channels cannot be JSON marshaled
// 	file, fileDone := tempFile(t)
// 	defer fileDone()

// 	err := fs.AppendFileJson(file.Name(), ch)
// 	assert.Error(t, err)
// }

// func Test_Touch(t *testing.T) {
// 	tempFile, fileDone := tempFileWithContent(t, "data")
// 	defer fileDone()

// 	tempDir, dirDone := tempDir(t)
// 	defer dirDone()

// 	filepath1 := filepath.Join(tempDir, "testfile.txt")
// 	err := fs.TouchFile(filepath1)
// 	failIf(t, err)
// 	assert.FileExists(t, filepath1)

// 	err = fs.TouchFile(tempFile.Name())
// 	failIf(t, err)
// 	assert.FileExists(t, tempFile.Name())

// 	content, err := fs.ReadFileString(tempFile.Name())
// 	failIf(t, err)
// 	assert.Equal(t, "data", content)

// 	err = fs.TouchFile("") // touching empty path should do nothing
// 	assert.Error(t, err)
// }

// func Test_EnsureFile(t *testing.T) {
// 	tempDir, dirDone := tempDir(t)
// 	defer dirDone()

// 	// is a directory = error
// 	filep := filepath.Join(tempDir, "subdir", "testfile.txt")
// 	err := fs.EnsureFile(tempDir)
// 	assert.ErrorAs(t, err, &fs.ErrIsDir)

// 	// create file
// 	err = fs.EnsureFile(filep)
// 	failIf(t, err)
// 	assert.True(t, fs.Exists(filep))
// 	assert.True(t, fs.IsFile(filep))

// 	// file exists = no error
// 	err = fs.EnsureFile(filep)
// 	failIf(t, err)
// 	assert.True(t, fs.Exists(filep))
// 	assert.True(t, fs.IsFile(filep))

// 	// invalid path
// 	err = fs.EnsureFile("")
// 	assert.Error(t, err)
// }

// func Test_ReplaceInFile(t *testing.T) {
// 	content := "Hello, World! World!"
// 	file, fileDone := tempFileWithContent(t, content)
// 	defer fileDone()

// 	err := fs.ReplaceInFile(file.Name(), []byte("World"), []byte("Go"))
// 	failIf(t, err)

// 	data, err := fs.ReadFileString(file.Name())
// 	failIf(t, err)
// 	assert.Equal(t, "Hello, Go! Go!", data)

// 	// old not found = no change
// 	err = fs.ReplaceInFile(file.Name(), []byte("Python"), []byte("Java"))
// 	failIf(t, err)

// 	data, err = fs.ReadFileString(file.Name())
// 	failIf(t, err)
// 	assert.Equal(t, "Hello, Go! Go!", data)

// 	// invalid file
// 	err = fs.ReplaceInFile(file.Name()+"-nonexistent", []byte("a"), []byte("b"))
// 	assert.Error(t, err)
// }

// func Test_ReplaceInFileString(t *testing.T) {
// 	content := "Hello, World! World!"
// 	file, fileDone := tempFileWithContent(t, content)
// 	defer fileDone()

// 	err := fs.ReplaceInFileString(file.Name(), "World", "Go")
// 	failIf(t, err)

// 	data, err := fs.ReadFileString(file.Name())
// 	failIf(t, err)
// 	assert.Equal(t, "Hello, Go! Go!", data)

// 	// old not found = no change
// 	err = fs.ReplaceInFileString(file.Name(), "Python", "Java")
// 	failIf(t, err)

// 	data, err = fs.ReadFileString(file.Name())
// 	failIf(t, err)
// 	assert.Equal(t, "Hello, Go! Go!", data)

// 	// invalid file
// 	err = fs.ReplaceInFileString(file.Name()+"-nonexistent", "a", "b")
// 	assert.Error(t, err)
// }

// func Test_TempFile(t *testing.T) {
// 	prefix := "mytempfile_"
// 	tempFilePath, err := fs.CreateTempFile(prefix)
// 	failIf(t, err)
// 	defer os.Remove(tempFilePath)

// 	assert.FileExists(t, tempFilePath)
// 	assert.Contains(t, filepath.Base(tempFilePath), prefix)

// 	// Create another temp file to ensure uniqueness
// 	tempFilePath2, err := fs.CreateTempFile(prefix)
// 	failIf(t, err)
// 	defer os.Remove(tempFilePath2)

// 	assert.FileExists(t, tempFilePath2)
// 	assert.Contains(t, filepath.Base(tempFilePath2), prefix)
// 	assert.NotEqual(t, tempFilePath, tempFilePath2)
// }

// func Test_TempFileOpen(t *testing.T) {
// 	prefix := "mytempfileopen_"
// 	tempFile, err := fs.CreateTempFileOpen(prefix)
// 	failIf(t, err)
// 	defer func() {
// 		tempFile.Close()
// 		os.Remove(tempFile.Name())
// 	}()

// 	assert.FileExists(t, tempFile.Name())
// 	assert.Contains(t, filepath.Base(tempFile.Name()), prefix)

// 	// Write to the temp file
// 	content := "Hello, TempFileOpen!"
// 	_, err = tempFile.WriteString(content)
// 	failIf(t, err)

// 	// Read back the content
// 	data, err := fs.ReadFileString(tempFile.Name())
// 	failIf(t, err)
// 	assert.Equal(t, content, data)
// }

// func Test_TempFile_InvalidPrefix(t *testing.T) {
// 	_, err := fs.CreateTempFile("")
// 	assert.NoError(t, err) // empty prefix is allowed

// 	// Extremely long prefix
// 	longPrefix := ""
// 	for i := 0; i < 300; i++ {
// 		longPrefix += "a"
// 	}
// 	_, err = fs.CreateTempFile(longPrefix)
// 	assert.Error(t, err) // long prefix is allowed, but may be truncated by the OS
// }

// func Test_IsSameFile(t *testing.T) {
// 	file1, file1Done := tempFile(t)
// 	defer file1Done()

// 	file2, file2Done := tempFile(t)
// 	defer file2Done()

// 	// Same file
// 	same := fs.IsSameFile(file1.Name(), file1.Name())
// 	assert.True(t, same)

// 	// Different files
// 	same = fs.IsSameFile(file1.Name(), file2.Name())
// 	assert.False(t, same)

// 	// Non-existent file
// 	same = fs.IsSameFile(file1.Name(), file2.Name()+"-nonexistent")
// 	assert.False(t, same)

// 	// Both non-existent files
// 	same = fs.IsSameFile(file1.Name()+"-nonexistent", file2.Name()+"-nonexistent")
// 	assert.False(t, same)
// }

// func Test_IsFileEmpty(t *testing.T) {
// 	file, fileDone := tempFile(t)
// 	defer fileDone()

// 	// Empty file
// 	empty, err := fs.IsEmptyFile(file.Name())
// 	failIf(t, err)
// 	assert.True(t, empty)

// 	// Non-empty file
// 	_, err = file.WriteString("data")
// 	failIf(t, err)
// 	file.Close() // close to flush

// 	empty, err = fs.IsEmptyFile(file.Name())
// 	failIf(t, err)
// 	assert.False(t, empty)

// 	// Non-existent file
// 	empty, err = fs.IsEmptyFile(file.Name() + "-nonexistent")
// 	assert.Error(t, err)
// 	assert.False(t, empty)
// }

// // DIRECTORIES

// func Test_IsDir(t *testing.T) {
// 	file, fileDone := tempFile(t)
// 	defer fileDone()
// 	assert.False(t, fs.IsDir(file.Name()))
// 	assert.False(t, fs.IsDir(file.Name()+"-nonexistent"))

// 	dir, dirDone := tempDir(t)
// 	defer dirDone()
// 	assert.True(t, fs.IsDir(dir))
// 	assert.False(t, fs.IsDir(dir+"-nonexistent"))
// }

// func Test_CreateDir(t *testing.T) {
// 	tempDir, dirDone := tempDir(t)
// 	defer dirDone()

// 	dirp := filepath.Join(tempDir, "subdir1", "subdir2")
// 	err := fs.CreateDir(dirp)
// 	failIf(t, err)
// 	assert.True(t, fs.Exists(dirp))
// 	assert.True(t, fs.IsDir(dirp))

// 	// Creating again should not cause error
// 	err = fs.CreateDir(dirp)
// 	failIf(t, err)
// 	assert.True(t, fs.Exists(dirp))
// 	assert.True(t, fs.IsDir(dirp))
// }

// func Test_EmptyDir(t *testing.T) {
// 	tempDir, dirDone := tempDir(t)
// 	defer dirDone()

// 	// is a file = error
// 	tempFile, fileDone := tempFile(t)
// 	defer fileDone()
// 	err := fs.EmptyDir(tempFile.Name())
// 	assert.ErrorAs(t, err, &fs.ErrNotDir)

// 	// non-existent dir = error
// 	nonexistentDir := filepath.Join(tempDir, "nonexistent")
// 	err = fs.EmptyDir(nonexistentDir)
// 	assert.ErrorAs(t, err, &fs.ErrNotDir)

// 	// create dir with files
// 	dirp := filepath.Join(tempDir, "subdir")
// 	err = fs.CreateDir(dirp)
// 	failIf(t, err)
// 	file1 := filepath.Join(dirp, "file1.txt")
// 	file2 := filepath.Join(dirp, "file2.txt")
// 	err = fs.WriteFileString(file1, "data1")
// 	failIf(t, err)
// 	err = fs.WriteFileString(file2, "data2")
// 	failIf(t, err)
// 	assert.True(t, fs.Exists(file1))
// 	assert.True(t, fs.Exists(file2))

// 	// empty dir
// 	err = fs.EmptyDir(dirp)
// 	failIf(t, err)
// 	assert.False(t, fs.Exists(file1))
// 	assert.False(t, fs.Exists(file2))
// 	assert.True(t, fs.Exists(dirp))
// 	assert.True(t, fs.IsDir(dirp))
// }

// func Test_IsEmptyDir(t *testing.T) {
// 	tempDir, dirDone := tempDir(t)
// 	defer dirDone()

// 	// is a file = error
// 	tempFile, fileDone := tempFile(t)
// 	defer fileDone()
// 	empty, err := fs.IsEmptyDir(tempFile.Name())
// 	assert.ErrorAs(t, err, &fs.ErrNotDir)
// 	assert.False(t, empty)

// 	// non-existent dir = error
// 	nonexistentDir := filepath.Join(tempDir, "nonexistent")
// 	empty, err = fs.IsEmptyDir(nonexistentDir)
// 	assert.ErrorAs(t, err, &fs.ErrNotDir)
// 	assert.False(t, empty)

// 	// empty dir
// 	dirp := filepath.Join(tempDir, "subdir")
// 	err = fs.CreateDir(dirp)
// 	failIf(t, err)
// 	empty, err = fs.IsEmptyDir(dirp)
// 	failIf(t, err)
// 	assert.True(t, empty)

// 	// dir with files
// 	file1 := filepath.Join(dirp, "file1.txt")
// 	err = fs.WriteFileString(file1, "data1")
// 	failIf(t, err)
// 	empty, err = fs.IsEmptyDir(dirp)
// 	failIf(t, err)
// 	assert.False(t, empty)
// }

// func Test_EnsureDir(t *testing.T) {
// 	tempDir, dirDone := tempDir(t)
// 	defer dirDone()

// 	tempFile, fileDone := tempFile(t)
// 	defer fileDone()

// 	err := fs.EnsureDir(tempFile.Name())
// 	assert.ErrorAs(t, err, &fs.ErrIsFile)

// 	dirp := filepath.Join(tempDir, "subdir1", "subdir2")
// 	err = fs.EnsureDir(dirp)
// 	failIf(t, err)
// 	assert.True(t, fs.Exists(dirp))
// 	assert.True(t, fs.IsDir(dirp))

// 	// Ensuring again should not cause error
// 	err = fs.EnsureDir(dirp)
// 	failIf(t, err)
// 	assert.True(t, fs.Exists(dirp))
// 	assert.True(t, fs.IsDir(dirp))
// }

// func Test_List(t *testing.T) {
// 	tempDir, dirDone := tempDir(t)
// 	defer dirDone()

// 	// is a file = error
// 	tempFile, fileDone := tempFile(t)
// 	defer fileDone()
// 	_, err := fs.List(tempFile.Name())
// 	assert.ErrorAs(t, err, &fs.ErrNotDir)

// 	// non-existent dir = error
// 	nonexistentDir := filepath.Join(tempDir, "nonexistent")
// 	_, err = fs.List(nonexistentDir)
// 	assert.ErrorAs(t, err, &fs.ErrNotDir)
// 	// empty dir
// 	dirp := filepath.Join(tempDir, "subdir")
// 	err = fs.CreateDir(dirp)
// 	failIf(t, err)
// 	entries, err := fs.List(dirp)
// 	failIf(t, err)
// 	assert.Equal(t, []string{}, entries)

// 	// dir with files and subdirs
// 	file1 := filepath.Join(dirp, "file1.txt")
// 	file2 := filepath.Join(dirp, "file2.txt")
// 	subdir := filepath.Join(dirp, "subdir2")
// 	fs.TouchFile(file1)
// 	fs.TouchFile(file2)
// 	err = fs.CreateDir(subdir)
// 	failIf(t, err)

// 	entries, err = fs.List(dirp)
// 	failIf(t, err)
// 	assert.ElementsMatch(t, []string{"file1.txt", "file2.txt", "subdir2"}, entries)
// }

// func Test_ListFiles(t *testing.T) {
// 	tempDir, dirDone := tempDir(t)
// 	defer dirDone()

// 	// is a file = error
// 	tempFile, fileDone := tempFile(t)
// 	defer fileDone()
// 	_, err := fs.ListFiles(tempFile.Name())
// 	assert.ErrorAs(t, err, &fs.ErrNotDir)

// 	// non-existent dir = error
// 	nonexistentDir := filepath.Join(tempDir, "nonexistent")
// 	_, err = fs.ListFiles(nonexistentDir)
// 	assert.ErrorAs(t, err, &fs.ErrNotDir)

// 	// empty dir
// 	dirp := filepath.Join(tempDir, "subdir")
// 	err = fs.CreateDir(dirp)
// 	failIf(t, err)
// 	files, err := fs.ListFiles(dirp)
// 	failIf(t, err)
// 	assert.Equal(t, []string{}, files)

// 	// dir with files and subdirs - should only return files
// 	file1 := filepath.Join(dirp, "file1.txt")
// 	file2 := filepath.Join(dirp, "file2.txt")
// 	subdir := filepath.Join(dirp, "subdir2")
// 	fs.TouchFile(file1)
// 	fs.TouchFile(file2)
// 	err = fs.CreateDir(subdir)
// 	failIf(t, err)

// 	files, err = fs.ListFiles(dirp)
// 	failIf(t, err)
// 	assert.ElementsMatch(t, []string{"file1.txt", "file2.txt"}, files)
// }

// func Test_ListDirs(t *testing.T) {
// 	tempDir, dirDone := tempDir(t)
// 	defer dirDone()

// 	// is a file = error
// 	tempFile, fileDone := tempFile(t)
// 	defer fileDone()
// 	_, err := fs.ListDirs(tempFile.Name())
// 	assert.ErrorAs(t, err, &fs.ErrNotDir)

// 	// non-existent dir = error
// 	nonexistentDir := filepath.Join(tempDir, "nonexistent")
// 	_, err = fs.ListDirs(nonexistentDir)
// 	assert.ErrorAs(t, err, &fs.ErrNotDir)

// 	// empty dir
// 	dirp := filepath.Join(tempDir, "subdir")
// 	err = fs.CreateDir(dirp)
// 	failIf(t, err)
// 	dirs, err := fs.ListDirs(dirp)
// 	failIf(t, err)
// 	assert.Equal(t, []string{}, dirs)

// 	// dir with files and subdirs - should only return directories
// 	file1 := filepath.Join(dirp, "file1.txt")
// 	file2 := filepath.Join(dirp, "file2.txt")
// 	subdir1 := filepath.Join(dirp, "subdir1")
// 	subdir2 := filepath.Join(dirp, "subdir2")
// 	fs.TouchFile(file1)
// 	fs.TouchFile(file2)
// 	err = fs.CreateDir(subdir1)
// 	failIf(t, err)
// 	err = fs.CreateDir(subdir2)
// 	failIf(t, err)

// 	dirs, err = fs.ListDirs(dirp)
// 	failIf(t, err)
// 	assert.ElementsMatch(t, []string{"subdir1", "subdir2"}, dirs)
// }

// func Test_ListRecursive(t *testing.T) {
// 	tempDir, dirDone := tempDir(t)
// 	defer dirDone()

// 	// is a file = error
// 	tempFile, fileDone := tempFile(t)
// 	defer fileDone()
// 	_, err := fs.ListRecursive(tempFile.Name())
// 	assert.ErrorAs(t, err, &fs.ErrNotDir)

// 	// non-existent dir = error
// 	nonexistentDir := filepath.Join(tempDir, "nonexistent")
// 	_, err = fs.ListRecursive(nonexistentDir)
// 	assert.ErrorAs(t, err, &fs.ErrNotDir)

// 	// empty dir
// 	dirp := filepath.Join(tempDir, "subdir")
// 	err = fs.CreateDir(dirp)
// 	failIf(t, err)
// 	entries, err := fs.ListRecursive(dirp)
// 	failIf(t, err)
// 	assert.Equal(t, []string{}, entries)

// 	// dir with nested files and subdirs
// 	file1 := filepath.Join(dirp, "file1.txt")
// 	subdir1 := filepath.Join(dirp, "subdir1")
// 	subdir2 := filepath.Join(subdir1, "subdir2")
// 	file2 := filepath.Join(subdir1, "file2.txt")
// 	file3 := filepath.Join(subdir2, "file3.txt")
// 	fs.TouchFile(file1)
// 	err = fs.CreateDir(subdir1)
// 	failIf(t, err)
// 	fs.TouchFile(file2)
// 	err = fs.CreateDir(subdir2)
// 	failIf(t, err)
// 	fs.TouchFile(file3)

// 	entries, err = fs.ListRecursive(dirp)
// 	failIf(t, err)
// 	expected := []string{
// 		"file1.txt",
// 		"subdir1",
// 		filepath.Join("subdir1", "file2.txt"),
// 		filepath.Join("subdir1", "subdir2"),
// 		filepath.Join("subdir1", "subdir2", "file3.txt"),
// 	}
// 	assert.ElementsMatch(t, expected, entries)
// }

// func Test_ListFilesRecursive(t *testing.T) {
// 	tempDir, dirDone := tempDir(t)
// 	defer dirDone()

// 	// is a file = error
// 	tempFile, fileDone := tempFile(t)
// 	defer fileDone()
// 	_, err := fs.ListFilesRecursive(tempFile.Name())
// 	assert.ErrorAs(t, err, &fs.ErrNotDir)

// 	// non-existent dir = error
// 	nonexistentDir := filepath.Join(tempDir, "nonexistent")
// 	_, err = fs.ListFilesRecursive(nonexistentDir)
// 	assert.ErrorAs(t, err, &fs.ErrNotDir)

// 	// empty dir
// 	dirp := filepath.Join(tempDir, "subdir")
// 	err = fs.CreateDir(dirp)
// 	failIf(t, err)
// 	files, err := fs.ListFilesRecursive(dirp)
// 	failIf(t, err)
// 	assert.Equal(t, []string{}, files)

// 	// dir with nested files and subdirs - should only return files
// 	file1 := filepath.Join(dirp, "file1.txt")
// 	subdir1 := filepath.Join(dirp, "subdir1")
// 	subdir2 := filepath.Join(subdir1, "subdir2")
// 	file2 := filepath.Join(subdir1, "file2.txt")
// 	file3 := filepath.Join(subdir2, "file3.txt")
// 	fs.TouchFile(file1)
// 	err = fs.CreateDir(subdir1)
// 	failIf(t, err)
// 	fs.TouchFile(file2)
// 	err = fs.CreateDir(subdir2)
// 	failIf(t, err)
// 	fs.TouchFile(file3)

// 	files, err = fs.ListFilesRecursive(dirp)
// 	failIf(t, err)
// 	expected := []string{
// 		"file1.txt",
// 		filepath.Join("subdir1", "file2.txt"),
// 		filepath.Join("subdir1", "subdir2", "file3.txt"),
// 	}
// 	assert.ElementsMatch(t, expected, files)
// }

// func Test_ListDirsRecursive(t *testing.T) {
// 	tempDir, dirDone := tempDir(t)
// 	defer dirDone()

// 	// is a file = error
// 	tempFile, fileDone := tempFile(t)
// 	defer fileDone()
// 	_, err := fs.ListDirsRecursive(tempFile.Name())
// 	assert.ErrorAs(t, err, &fs.ErrNotDir)

// 	// non-existent dir = error
// 	nonexistentDir := filepath.Join(tempDir, "nonexistent")
// 	_, err = fs.ListDirsRecursive(nonexistentDir)
// 	assert.ErrorAs(t, err, &fs.ErrNotDir)

// 	// empty dir
// 	dirp := filepath.Join(tempDir, "subdir")
// 	err = fs.CreateDir(dirp)
// 	failIf(t, err)
// 	dirs, err := fs.ListDirsRecursive(dirp)
// 	failIf(t, err)
// 	assert.Equal(t, []string{}, dirs)

// 	// dir with nested files and subdirs - should only return directories
// 	file1 := filepath.Join(dirp, "file1.txt")
// 	subdir1 := filepath.Join(dirp, "subdir1")
// 	subdir2 := filepath.Join(subdir1, "subdir2")
// 	file2 := filepath.Join(subdir1, "file2.txt")
// 	file3 := filepath.Join(subdir2, "file3.txt")
// 	fs.TouchFile(file1)
// 	err = fs.CreateDir(subdir1)
// 	failIf(t, err)
// 	fs.TouchFile(file2)
// 	err = fs.CreateDir(subdir2)
// 	failIf(t, err)
// 	fs.TouchFile(file3)

// 	dirs, err = fs.ListDirsRecursive(dirp)
// 	failIf(t, err)
// 	expected := []string{
// 		"subdir1",
// 		filepath.Join("subdir1", "subdir2"),
// 	}
// 	assert.ElementsMatch(t, expected, dirs)
// }

// func Test_Glob(t *testing.T) {
// 	tempDir, dirDone := tempDir(t)
// 	defer dirDone()

// 	// empty dir
// 	dirp := filepath.Join(tempDir, "subdir")
// 	err := fs.CreateDir(dirp)
// 	failIf(t, err)
// 	matches, err := fs.Glob(dirp, "*.txt")
// 	failIf(t, err)
// 	assert.Equal(t, []string{}, matches)

// 	// dir with files and subdirs
// 	file1 := filepath.Join(dirp, "file1.txt")
// 	file2 := filepath.Join(dirp, "file2.log")
// 	subdir := filepath.Join(dirp, "subdir")
// 	fs.TouchFile(file1)
// 	fs.TouchFile(file2)
// 	err = fs.CreateDir(subdir)
// 	failIf(t, err)

// 	matches, err = fs.Glob(dirp, "*.txt")
// 	failIf(t, err)
// 	assert.Equal(t, []string{"file1.txt"}, matches)

// 	matches, err = fs.Glob(dirp, "*.log")
// 	failIf(t, err)
// 	assert.Equal(t, []string{"file2.log"}, matches)

// 	matches, err = fs.Glob(dirp, "*.*")
// 	failIf(t, err)
// 	assert.ElementsMatch(t, []string{"file1.txt", "file2.log"}, matches)

// 	matches, err = fs.Glob(dirp, "*.md")
// 	failIf(t, err)
// 	assert.Equal(t, []string{}, matches)
// }

// func Test_CreateTempDir(t *testing.T) {
// 	prefix := "mytempdir_"
// 	tempDirPath, err := fs.CreateTempDir(prefix)
// 	failIf(t, err)
// 	defer os.RemoveAll(tempDirPath)

// 	assert.DirExists(t, tempDirPath)
// 	assert.Contains(t, filepath.Base(tempDirPath), prefix)

// 	// Create another temp dir to ensure uniqueness
// 	tempDirPath2, err := fs.CreateTempDir(prefix)
// 	failIf(t, err)
// 	defer os.RemoveAll(tempDirPath2)

// 	assert.DirExists(t, tempDirPath2)
// 	assert.Contains(t, filepath.Base(tempDirPath2), prefix)
// 	assert.NotEqual(t, tempDirPath, tempDirPath2)
// }

// // GENERAL

// func Test_Exists(t *testing.T) {
// 	file, fileDone := tempFile(t)
// 	defer fileDone()
// 	assert.True(t, fs.Exists(file.Name()))
// 	assert.False(t, fs.Exists(file.Name()+"-nonexistent"))

// 	dir, dirDone := tempDir(t)
// 	defer dirDone()
// 	assert.True(t, fs.Exists(dir))
// 	assert.False(t, fs.Exists(dir+"-nonexistent"))
// }

// // PATH

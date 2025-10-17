# File System Helpers

This package tries to provide a more sane and cohesive way to interact with the file system.

<!-- TOC -->

- [File System Helpers](#file-system-helpers)
  - [Installing](#installing)
  - [Basic Usage](#basic-usage)
  - [Naming Convention](#naming-convention)
  - [Path Anatomy](#path-anatomy)
  - [Watching](#watching)

<!-- /TOC -->

## Installing

```
go get github.com/renatopp/go-fs
```

## Basic Usage

```
import "github.com/renatopp/go-fs"

func main() {
  if fs.IsDir(".") {
    println("ok")
  }
}
```

## Naming Convention

| Function | Description | Examples | 
| --- | --- | --- |
| `*File` | affects exclusively files or returns files, probably resulting in an error if a directory is provided. | `ReadFile()` |
| `*Dir` | affects exclusively directories or returns directories. | `EmptyDir()` `GetHomeDir()` | 
| `*Path` | affects the path string (not the filesystem, just the string itself). | `JoinPath()` `GetPathName()` |
| `Force*` | ignore errors and return just the values. Only applicable for functions that return (value, error). | `ForceReadFile()` |
| `Other` | accept directories or files, handling them differently if necessary. | `Hide()` |

## Path Anatomy

You can use `GetPathParts` to extract all these infos at the same time.

**Absolute**

It is the complete and absolute path, including every part.

```
/c/users/dev/fs/path.go
^^^^^^^^^^^^^^^^^^^^^^^
       Absolute
```

**Base**

It is the name of the file (or folder) including its extension.

```
/c/users/dev/fs/path.go
                ^^^^^^^
                 Base
```

**Name**

Represents the name of the file (or folder) excluding its extension.

```
/c/users/dev/fs/path.go
                ^^^^
                Name
```

**Extension**

It is the extension of the file, including the dot.

```
/c/users/dev/fs/path.go
                    ^^^
                    Ext
```

**Extension Name**

Same as the extension, but excluding the dot.

```
/c/users/dev/fs/path.go
                     ^^
                   ExtName
```

**Parent**

Called `Dir` in other libraries, the Parent is the path excluding the last part, usually the file or last directory name.

```
/c/users/dev/fs/path.go
^^^^^^^^^^^^^^^
     Parent
```

**Parent Name**

Is is the name of the last part before the file or last directory name. Equivalent to `BasePath(ParentPath(path))`.

```
/c/users/dev/fs/path.go
             ^^
         ParentName
```

**Volume**

For windows only, represents the volume of the path.

```
/c/users/dev/fs/path.go
 ^
 Volume
```

## Watching

This library uses [fsnotify](https://github.com/fsnotify/fsnotify) internally to watch for directory and files changes.

```go
fs.Watch(context.Background(), "sample/", func (event fs.Event) {
  println(event.Path, "has changed")
})
```

The event have the following information:

- `Op`, a bitmasked int describing the event (or events) that happened. You can check it for specific event as `event.Op.Has(fs.EvtCreate)`:
  -	EvtCreate
  -	EvtRemove
  -	EvtWrite
  -	EvtRename
  -	EvtChmod
  -	EvtError
- `Path`, the path of the file or folder that generated the event. It will be prefixed by the watched path provided (`sample/file.txt` in the example above).
- `Err`, for error events.

Other options:

```go
fs.NewWatcher().Watch(...)
fs.Watch(...)
fs.WatchRecursive(...)
fs.WatchGlob(...)
```



package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/renatopp/go-fs"
)

type Writer struct {
	w *tabwriter.Writer
}

func (w *Writer) Write(s ...any) {
	for _, str := range s {
		w.w.Write([]byte(fmt.Sprintf("%v", str) + "\t"))
	}
	w.w.Write([]byte("\n"))
}

func (w *Writer) Flush() {
	w.w.Flush()
}

func main() {
	writer := &Writer{w: tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)}
	defer writer.Flush()

	//..
	println("Watching...") //
	fs.Watch(context.Background(), "example", func(event fs.Event) {
		fmt.Printf("Event: %s, Path: %s, Err: %v\n", event.String(), event.Path, event.Err)
	})

	// sep := ""
	// writer.Write(sep, sep)
	// writer.Write("File System Utilities")
	// writer.Write(sep, sep)
	// writer.Write("fs.PathSeparator", fs.PathSeparator())
	// writer.Write(sep, sep)
	// writer.Write(sep, sep)
	// writer.Write("fs.CurrentPath()", fs.Force(fs.CurrentPath()))
	// writer.Write("fs.HomePath()", fs.Force(fs.HomePath()))
	// writer.Write("fs.CachePath()", fs.Force(fs.CachePath()))
	// writer.Write("fs.TempPath()", fs.TempPath())
	// writer.Write("fs.ConfigPath()", fs.Force(fs.ConfigPath()))
	// writer.Write(sep, sep)
	// writer.Write("function", ".", "..", "~", "main.go", "example", "nonexistent")
	// writer.Write(sep, sep)
	// writer.Write("AbsolutePath",
	// 	fs.Force(fs.AbsolutePath(".")),
	// 	fs.Force(fs.AbsolutePath("..")),
	// 	fs.Force(fs.AbsolutePath("~")),
	// 	fs.Force(fs.AbsolutePath("main.go")),
	// 	fs.Force(fs.AbsolutePath("example")),
	// 	fs.Force(fs.AbsolutePath("nonexistent")),
	// )
	// writer.Write("fs.AbsolutePath(..)", fs.Force(fs.AbsolutePath("..")))
	// writer.Write("fs.AbsolutePath(~)", fs.Force(fs.AbsolutePath("~")))
	// writer.Write("fs.RelativePath(.., .)", fs.Force(fs.RelativePath("..", ".")))
	// writer.Write("fs.ExtensionPath(.)", fs.ExtensionPath("."))
	// writer.Write("fs.ExtensionPath(.txt)", fs.ExtensionPath(".txt"))
	// writer.Write("fs.ExtensionPath(name.txt)", fs.ExtensionPath("name.txt"))
	// writer.Write("fs.IsHidden(name.txt)", fs.Force(fs.IsHidden(".")))
	// writer.Write("fs.IsHidden(.txt)", fs.Force(fs.IsHidden(".")))
	// fs.IsFile2 = internal.IsFile
}

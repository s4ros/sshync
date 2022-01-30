package main

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

func _error(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	userHomeDirectory, err := os.UserHomeDir()
	_error(err)
	src := filepath.Join(userHomeDirectory, ".ssh")
	t := time.Now().UnixMilli()
	out := fmt.Sprint("ssh-archive-", t, ".tar")

	filenames := getAllFiles(src)

	err = createArchive(out, filenames)
	_error(err)
}

func getAllFiles(root string) []string {
	filenames := []string{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			// skip directories
			return nil
		}
		filenames = append(filenames, path)
		return nil
	})
	_error(err)
	return filenames
}

func createArchive(out string, filenames []string) error {
	fd_out, err := os.Create(out)
	_error(err)
	defer fd_out.Close()

	tw := tar.NewWriter(fd_out)
	defer tw.Close()

	for _, file := range filenames {
		err := addToArchive(tw, file)
		_error(err)
	}
	return nil
}

func addToArchive(tw *tar.Writer, filename string) error {
	fd, err := os.Open(filename)
	_error(err)
	defer fd.Close()

	fstat, err := fd.Stat()
	_error(err)

	header, err := tar.FileInfoHeader(fstat, fstat.Name())
	_error(err)

	header.Name = filename

	err = tw.WriteHeader(header)
	_error(err)

	fmt.Println("Archiving:", filename)
	_, err = io.Copy(tw, fd)
	_error(err)

	return nil
}

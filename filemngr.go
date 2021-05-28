package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type FileManager struct {
	Path    string
	ch      chan string
	running bool
}

func (fm FileManager) WriteChan() chan string {
	return fm.ch
}

func (fm FileManager) Tail(w io.Writer, lines int) error {
	idx := 0
	ll := make([]string, lines)
	hasContent := false

	f, err := os.Open(fm.Path)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		ll[idx] = scanner.Text()
		idx++

		if idx >= lines {
			idx = 0
		}

		hasContent = true
	}

	if !hasContent {
		return nil
	}

	//TODO put log when write fails
	//write in order
	for i := idx; i < lines; i++ {
		if ll[i] == "" {
			continue
		}

		fmt.Fprintf(w, "%s\n", ll[i])
	}

	//jump to 0 index
	for i := 0; i < idx; i++ {
		if ll[i] == "" {
			break
		}

		fmt.Fprintf(w, "%s\n", ll[i])
	}

	return nil
}

func (fm FileManager) Run() error {
	if fm.running {
		return errors.New("Already running")
	}

	fm.running = true
	file := openOrCreateFile(fm.Path)
	defer closeFile(file)
	for {
		select {
		case t := <-fm.ch:
			if file == nil {
				file = openOrCreateFile(fm.Path)
			}

			if file != nil {
				_, err := file.WriteString(t)
				if err != nil {
					fmt.Println(err.Error())
				}

			}

		default:
			closeFile(file)
			file = nil
			time.Sleep(5 * time.Second)
		}
	}
}

func NewFileManager(path string) *FileManager {
	return &FileManager{path, make(chan string, 100), false}
}

func closeFile(f *os.File) {
	if f != nil {
		f.Close()
	}
}

func openOrCreateFile(p string) *os.File {
	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}

	return f
}

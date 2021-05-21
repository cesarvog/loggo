package main

import (
		"time"
		"fmt"
		"bufio"
		"os"
		"strings"
		"errors"
)

type FileManager struct {
	Path string
	ch chan string
	running bool
}

func (fm FileManager) WriteChan() chan string {
	return fm.ch
}

func (fm FileManager) Tail(lines int) (string, error) {
	idx := 0
	ll := make([]string, lines)
	hasContent := false

	f, err := os.Open(fm.Path)
	if err != nil {
		return "", nil
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		ll[idx] = scanner.Text()
		idx++

		if(idx >= lines) {
			idx = 0
		}

		hasContent = true
	}

	if !hasContent {
		return "", nil
	}

	//put array in order
	lo := make([]string, lines)
	j := 0

	for i := idx; i < lines; i++ {
		if ll[i] == "" {
			continue
		}

		lo[j] = ll[i]
		j++
	}

	//jump to 0 index
	for i := 0; i<idx; i++ {
		if ll[i] == "" {
			break
		}

		lo[j] = ll[i]
		j++
	}

	return strings.Join(lo, ""), nil
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

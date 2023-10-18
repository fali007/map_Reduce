package search

import (
	"io"
	"os"
	"fmt"
	"time"
	"bufio"
	"strings"
)

var o chan string

func set_channel() {
	if o == nil {
		o = make(chan string)
	}
}

func get_files(folder string) []string {
	files, err := os.ReadDir(folder)
	if err != nil {
		panic(err)
	}
	s := []string{}
	for _,file := range files {
		if file.IsDir() {
			s = append(s, get_files(fmt.Sprintf("%s/%s", folder, file.Name()))...)
			continue
		} else{
			s = append(s, fmt.Sprintf("%s/%s", folder, file.Name()))
		}
	}
	return s
}

func read_chunks (f *os.File, key, filename string) {
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if strings.Contains(string(line), key) {
			o <- filename
			break
		}
	}
}

func read_active () {
	for e := range o {
		fmt.Println(e)
	}
}

func Search (key, folder string) {
	set_channel()
	go read_active()

	files := get_files(folder)

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		go read_chunks(f, key, file)
	}
	time.Sleep(1 * time.Second)
	close(o)
}
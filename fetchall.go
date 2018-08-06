package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)
	count := 0
	for _, url := range os.Args[1:] {
		count++
		name := "file" + strconv.Itoa(count)
		go fetch(url, ch, name) // start a gorutine
	}
	for range os.Args[1:] {
		fmt.Println(<-ch) // receive from channel ch
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string, file_name string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) // send to channel ch
		return
	}

	// Creates a file to write the content
	file, err := os.Create(file_name)
	if err != nil {
		ch <- fmt.Sprintf("while creating file  %s: %v", file_name, err)
		return
	}

	// Copy the content of the file and get the size in nbytes
	nbytes, err := io.Copy(file, resp.Body)
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}

	resp.Body.Close() // don't leak resources
	file.Close()

	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s %s", secs, nbytes, file_name, url)
}

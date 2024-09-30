package main

import(
	"fmt"
	"log"
//	"net/http"
	"os"
//	"io"
	"bufio"
	"sync"
	"time"
	"strings"
)

func readData(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := bufio.NewReader(f) // f is io.Reader, but not buffered
	scanner := bufio.NewScanner(r)
	// type SplitFunc func(data []byte, atEOF bool) \
	//   (advance int, token []byte, err error)
	scanner.Split(bufio.ScanLines)  // Assigning the function signature
	data := make([]string, 0)
	for scanner.Scan() {
		t := scanner.Text()
		if strings.Trim(t) == "" {
			continue
		}
		data = append(data, scanner.Text())
	}
	return data
}

func checkMonUrl(stop chan struct{}, url string) {
	for {
		<-stop
		fmt.Println("Checking monitored URL", url)
		time.Sleep(time.Second)
	}
	return
}

func checkRefUrl(wg *sync.WaitGroup, url string) {
	defer wg.Done()
	fmt.Println("Checking reference URL", url)
	time.Sleep(2*time.Second)
	return
}

func main() {
	refUrls := readData("./reference-urls.txt")
	monUrls := readData("./monitored-urls.txt")
	wg := &sync.WaitGroup{}
	// wg.Add(len(refUrls) + len(monUrls))
	wg.Add(len(refUrls))
	stopper := make(chan struct{})
	for _, ru := range refUrls {
		go checkMonUrl(stopper, ru)
	}
	for _, mu := range monUrls {
		go checkRefUrl(wg, mu)
	}
	wg.Wait()
	time.Sleep(10*time.Second)
	close(stopper) // Closing channel will stop "go checkMonUrl(...)"
}
		

/*
	r, err := http.Get("http://www.google.com/")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	for k, v := range(r.Header) {
		fmt.Printf("%s: %s\n", k, v)
	}
	fmt.Println("--------- Content ----------")

	part := make([]byte, 0xF)
	n, err2 := r.Body.Read(part)
	if err2 == io.EOF {
		defer fmt.Println("Got the whole body until io.EOF")
	} else if err2 != nil {
		log.Fatal("Error occurred", err2)
	}
	fmt.Println(string(part))
	fmt.Println(n, "bytes read in total")
*/



package main

import(
	"fmt"
	"log"
//	"net/http"
	"os"
//	"io"
	"bufio"
	"time"
	"strings"
	"sync"
//	"math/rand"
)

func readData(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	scanner := bufio.NewScanner(r)  // f was also io.Reader, but not buffered
	scanner.Split(bufio.ScanLines)  // Assigning the function signature
	data := make([]string, 0)
	for scanner.Scan() {  // just advance and give ok
		t := scanner.Text()
		if strings.TrimSpace(t) == "" {
			continue
		}
		data = append(data, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return data
}

func checkMonUrl(wg *sync.WaitGroup, reporter chan string, url string) {
	defer wg.Done()
	fmt.Println("In checkMonUrl()")
	fmt.Println("Writing for url", url)
	// To do: Probe here and report
	reporter <- fmt.Sprintf("Let's say %s is up", url)
	// time.Sleep(time.Duration(rand.Intn(3))*time.Second)
	return
}

func schedule(c *bool) {
	// Anything that sets c to false when monitoring should stop
	// and leave true to continue for another round:
	time.Sleep(5*time.Second)
	*c = false
}

func main() {
	monUrls := readData("./monitored-urls.txt")
	ch := make(chan string, len(monUrls))  // To do: Make buffer dynamic
	var wg sync.WaitGroup

	cont := true
	go schedule(&cont)

	for {
		fmt.Println("Check if we continue.")
		if (! cont) {
			fmt.Println("  Stop!")
			break
		} else {
			fmt.Println("  Continue!")
		}
		for _, url := range(monUrls) {
			wg.Add(1)
			fmt.Println("Starting a monitor")
			go checkMonUrl(&wg, ch, url)
		}
		wg.Wait() // Perhaps not needed
		for i:=0; i<len(monUrls); i++ {
			fmt.Println(<-ch)
		}
		// time.Sleep(time.Second)
	}
	fmt.Println("Before closing reporter")
	close(ch)
	fmt.Println("Leaving main()")
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



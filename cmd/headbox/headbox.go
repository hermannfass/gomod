package main

import (
	"log"
	"os"
	"io/ioutil"  // Read files
	"bufio"
	"fmt"
	"strings"
	"github.com/hermannfass/gomod/textboxes"
)

func main() {
	// Reader for user input
	reader := bufio.NewReader(os.Stdin) // Stdin implements io.Reader interface

	var style string
	styles := map[string]string{
		"a": "ascii",
		"s": "single",
		"d": "double",
		"m": "mixed",
		"": "mixed", // default
	}
	for style = ""; style == ""; {
		fmt.Print("Style (a)scii, (s)ingle, (d)ouble, (m)ixed? [m] ")
		s, err := reader.ReadString('\n')
		if err != nil { log.Fatal(err) }
		style = styles[strings.TrimSpace(s)]
	}

	var bw int = 76
	fmt.Print("Box Width [76]: ")
	input, err := reader.ReadString('\n')
	if err != nil { log.Fatal(err); }
	bws := strings.TrimSpace(input)
	if bws != "" {
		_, err := fmt.Sscan(bws, &bw)  // Skipping result count
		if (err != nil) {
			fmt.Printf("Input is not a number. Using default: %d\n", bw)
		}
	}

	texts := make(map[string]string, len(textboxes.TextFields))
	for _, n := range textboxes.TextFields {
		fmt.Printf("%s: ", n)
		input, err := reader.ReadString('\n')
		if err != nil { log.Fatal(err) }
		texts[n] = strings.TrimSpace(input)
	}

	var inPath string
	if len(os.Args) > 1 {
		inPath = os.Args[1]
		fmt.Println("Prepending header to file", inPath)
	} else {
		fmt.Print("Input File [none]: ")
		input, err := reader.ReadString('\n')
		if err != nil { log.Fatal(err); }
		inPath = strings.TrimSpace(input)
	}

	fmt.Print("Output File [console]: ")
	path, err := reader.ReadString('\n')
	if err != nil { log.Fatal(err); }
	outPath := strings.TrimSpace(path)
	w := bufio.NewWriter(os.Stdout)
	if outPath != "" {
		f, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		w = bufio.NewWriter(f)
	} 
	defer w.Flush()

	header := textboxes.HeaderBox(bw, style, texts)

	// Read existing file if specified and existing
	var content string // empty string ""
	if inPath != "" {
		if _, rErr := os.Stat(inPath); !os.IsNotExist(rErr) {
			data, rErr := ioutil.ReadFile(inPath)
			if (rErr != nil) { log.Fatal(rErr) }
			content = string(data)
		} else {
			fmt.Printf("File %s not existing. Input skipped.\n", inPath)
		}
	}

	output := header + content

	_, wErr := w.WriteString(output) // Not using number of bytes written
	if wErr != nil {
		log.Fatal(err)
	}

}


package main

import (
	"bufio"
	"os"
	"fmt"
	"log"
	"strings"
	"io/ioutil"  // Read files
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
		if err != nil { log.Fatal(err); return }
		inPath = strings.TrimSpace(input)
	}

	fmt.Print("Output File [" + inPath + "]: ")
	path, err := reader.ReadString('\n')
	if err != nil { log.Fatal(err); return }
	outPath := strings.TrimSpace(path)
	if outPath == "" {
		outPath = inPath
	} 

	header := textboxes.HeaderBox(style, texts)
	var content string // empty string ""
	if _, rErr := os.Stat(inPath); !os.IsNotExist(rErr) {
		data, rErr := ioutil.ReadFile(inPath)
		if (rErr != nil) { log.Fatal(rErr) }
		content = string(data)
	} else {
		fmt.Printf("File %s not existing. Using empty file.\n", inPath)
	}

	output := header + content
	wErr := ioutil.WriteFile(outPath, []byte(output), 0644)
	if wErr != nil {
		fmt.Println("Problem writing content to", outPath)
		log.Fatal(wErr);
		return
	}

}


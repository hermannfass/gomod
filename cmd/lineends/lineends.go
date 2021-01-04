package main

import (
	"fmt"
	"flag"
	"strings"
	"github.com/hermannfass/gomod/lineends"
	"log"
	"regexp"
)

func main() {
	/* lineends [-s targetformat] [-i] inputPath [-o] outputPath
	   lineends new.txt old.txt (anything to unix format)
	   lineends -o new.txt old.txt
	   lineends -i old.txt -o new.txt
	*/
	var sp, ip, op *string // pointers to system, input/output path
	var inFilePath, outFilePath string
	var data string

	sp = flag.String("s", "unix", "Target format: unix, dos, classicmac")
	ip = flag.String("i", "", "Input file (optional if not STDIN)") 
	op = flag.String("o", "", "Output file (optional if not STDOUT)")
	flag.Parse()
	argsLeft := flag.Args()
	if *ip != "" {
		inFilePath = *ip
	} else if len(argsLeft) > 0 {
		inFilePath, argsLeft = argsLeft[0], argsLeft[1:]
	} else {
		// No -i flag and no unflagged arg =>
		// To do: Help on usage: no -i and no unflagged argument
	}
	if *op != "" {
		outFilePath = *op
	} else if len(argsLeft) > 0 {
		outFilePath = argsLeft[0]
	}

	s := strings.ToLower(*sp)
	var le string = lineends.LineEnd[s]
	
	data = lineends.ReadData(inFilePath)

	re, err := regexp.Compile(`\r\n|\r|\n`)
	if err != nil {
		log.Fatal(err)
	}
	leCounter := make(map[string]int)
	leMatches := re.FindAllString(data, -1)
	for _, v := range(leMatches) {
		leCounter[v]++  // Before defined it has zero value, thus ok.
	}
	for k, v := range(leCounter) {
		fmt.Printf("%d matches for %#v\n", v, k)
	}

	data = re.ReplaceAllString(data, le)
	lineends.WriteData(outFilePath, data)
	
}


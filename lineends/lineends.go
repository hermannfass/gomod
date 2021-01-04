package lineends

import (
	"log"
	"os"
	"io"
	"bufio"
)

// Read data from a file (path to file as argument)
// or from a pipe (empty string as argument).
// To do: Could also use a URL to read from instead.
func ReadData(fn string) string {
	var data []rune  // result as read (runes)
	r := bufio.NewReader(os.Stdin) // Read from pipe 
	if (fn != "") {  // To do: Could also check for URLs and read via HTTP
		f, err := os.Open(fn) // Open for reading
		if (err != nil) {
			log.Fatal(err)
		}
		defer f.Close()
		r = bufio.NewReader(f)
	}
	for {
		in, _, err := r.ReadRune() // Not using width (byte count) of rune
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)  // calls os.Exit(1)
		}
		data = append(data, in)
	}
	return string(data)
}

func WriteData(fn string, data string) {
	w := bufio.NewWriter(os.Stdout)
	if (fn != "") {
		f, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		w = bufio.NewWriter(f)
	}
	defer w.Flush()
	_, err := w.WriteString(data) // Not using number of bytes written
	if err != nil {
		log.Fatal(err)
	}
}


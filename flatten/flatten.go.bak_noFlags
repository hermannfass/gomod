package main

import(
	"fmt"
	"flag"
	"log"
	"os"
	"path/filepath"
	"io"
	// "strings"
)

// Separator between individual path elements of filenames:
const pathSep string = "--"  // To do: This should be configurable

func main() {
	flag.Usage = func() {
		fmt.Println(usageText())
	} 
	flag.Parse()
	topdir := flag.Arg(0)
	var outdir string
	if flag.NArg() > 1 {
		outdir = flag.Arg(1)
	} else {
		outdir = topdir
	}
	fmt.Println("Args count: ", flag.NArg())
	fmt.Println("Looking into directory", topdir)
	flatten(topdir, pathToPrefix(topdir), outdir)
}


func flatten(dir string, prefix string, outdir string) {
	fmt.Println("Flattening", dir, "Prefix:", prefix)
	entries, _ := os.ReadDir(dir)  // to do: handle potential err
	for _, fi := range entries {
		if fi.IsDir() {
			newpre := prefix + pathSep + fi.Name()
			from := filepath.Join(dir, fi.Name()) 
			flatten(from, newpre, outdir)
		} else {
			// new file name with directory information included
			newfn := prefix + pathSep + fi.Name()
			// full path of the new (renamed) file:
			fmt.Println("outdir:", outdir)
			oldpath := filepath.Join(dir, fi.Name())
			newpath := filepath.Join(outdir, newfn)
			fmt.Println("prefix: ", prefix)
			fmt.Println("newpath: ", newpath)
			// fmt.Println(fi.Name(), "to move to:", newpath)
			fmt.Println("copy", oldpath, "to", newpath)
			filecopy(oldpath, newpath)
		}
	}
}

func filecopy(from, to string) {
	fdFrom, err := os.Open(from)
	if err != nil {
		log.Fatal(err)
	}
	defer fdFrom.Close()
	fdTo, err := os.Create(to)
	if err != nil {
		log.Fatal(err)
	}
	defer fdTo.Close()
	bytes, err := io.Copy(fdTo, fdFrom)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Copied %d bytes from %s to %s\n", bytes, from, to)
}
	
	

func pathToPrefix(p string) string {
	r := filepath.Base(p)
	fmt.Println("pathToPrefix returns:", r)
	return(r)
}


func usageText() string {
	return `
PURPOSE
Tool for copying files from a directory tree to one single directory
and prepending the directory names into the individual file names.

CALL
   flatten <top-directory> <target-directory>

`
}


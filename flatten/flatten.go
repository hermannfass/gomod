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
// const pathSep string = "--"  // To do: This should be configurable
var pathSep string // Added later; it might also become a pointer.
var v *bool        // Flag switching to verbose mode

func main() {
	flag.Usage = func() {
		fmt.Println(usageText())
		flag.PrintDefaults()
	} 
	pathSepFlag := flag.String("s", "--", "Path separator in output filenames")
	v = flag.Bool("v", false, "Verbose messages during execution")
	flag.Parse()
	pathSep = *pathSepFlag
	topdir := flag.Arg(0)
	var outdir string
	if flag.NArg() > 1 {
		outdir = flag.Arg(1)
		if _, fcerr := os.Stat(outdir); os.IsNotExist(fcerr) {
			// To do: Might ask if it should get created.
			fmt.Printf("The output directory %s does not exist\n", outdir)
			fmt.Println("Create it or select an existing directory.")
			os.Exit(1)
		}
	} else {
		fmt.Println("No output directory specified. Writing to", topdir)
		// To do: Could ask if it is ok to use the start directory.
		outdir = topdir
	}
	// fmt.Println("Args count: ", flag.NArg())
	flatten(topdir, pathToPrefix(topdir), outdir)
}


func flatten(dir string, prefix string, outdir string) {
	if *v {
		fmt.Printf("Directory %s\n  => Prefix %s\n", dir, prefix)
	}
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
			oldpath := filepath.Join(dir, fi.Name())
			newpath := filepath.Join(outdir, newfn)
			// fmt.Println("prefix: ", prefix)
			// fmt.Println("newpath: ", newpath)
			if *v {
				fmt.Printf("  Copying file %s\n", fi.Name())
				fmt.Printf("    to %s\n", newpath)
			}
			// fmt.Println("copy", oldpath, "to", newpath)
			b := filecopy(oldpath, newpath)
			if *v { fmt.Printf("    %d bytes written\n", b) }
		}
	}
}

func filecopy(from, to string) int64 {
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
	return(bytes)
}
	
	

func pathToPrefix(p string) string {
	r := filepath.Base(p)
	// fmt.Println("pathToPrefix returns:", r)
	return(r)
}


func usageText() string {
	return `
PURPOSE
Tool for copying files from a directory tree to one single directory
and prepending the directory names into the individual file names.

CALL
   flatten [-s <filename-separator>] <top-directory> <target-directory>

`
}


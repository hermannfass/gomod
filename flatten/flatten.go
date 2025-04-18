package main

import(
	"fmt"
	"flag" 
	"log"
	"os"
	"path/filepath"
	"io"
	"io/fs"
	"strings"
	"errors"
)

func main() {
	flag.Usage = func() {
		fmt.Println(usageText()) // No extra flag.PrintDefaults()
	}
	pathSepPtr := flag.String("sep", "_", "Path separator in output filenames")
	flag.Parse()
	pathSep := *pathSepPtr

	topdir, err := getTopDir(flag.Args())
	if err != nil {
		fmt.Println("Aborting.", err)
		os.Exit(0)
	} else if topdir == "" {
		fmt.Println("No directory to flatten, thus stopping here.")
		os.Exit(0)
	}

	outdir, err := getOutDir(topdir, flag.Args())
	if err != nil {
		fmt.Printf("Could not determine output directory. %s\n", err)
		os.Exit(0)
	} else if outdir == "" {
		fmt.Println("No output directory specified/accepted. Aborting.")
		os.Exit(0)
	}
	fmt.Println("Equivalent start prefix:     ", pathToPrefix(topdir))
	fmt.Println("Output directory:            ", outdir)
	fmt.Println("Path separator in new files: ", pathSep)
	fmt.Print(  "Is it ok to go ahead?  ")
	if (! answer()) {
		fmt.Println("=> Exiting")
		os.Exit(0)
	}
	flatten(topdir, pathToPrefix(topdir), outdir, pathSep)
	fmt.Println("Done")
}

func getTopDir(args []string) (string, error) {
	if len(args) < 1 {
		fmt.Println("No directory (argument 1) specified.")
		fmt.Println("Checking the current working directory.")
		r, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("Error determining working directory. %w", err)
		} else {
			fmt.Printf("Working directory is: %q.\n", r)
			fmt.Print("Do you want to flatten that?")
			if answer() {
				return r, nil
			} else {
				return "", nil
			}
		}
	} else {
		ok, err := existsAndIsDir(args[0])
		if err != nil {
			return "", err
		} else if ok {
			return args[0], nil
		} else {
			return "", nil
		}
	}
}


func getOutDir(topdir string, args []string) (string, error)  {
	if len(args) < 2 {
		fmt.Println("No output directory specified.")
		fmt.Printf("Should the output go to %q?", topdir)
		if answer() {
			return topdir, nil
		} else {
			return "", nil
		}
	} else {
		ok, err := existsAndIsDir(args[1])
		if ok {
			return args[1], nil
		} else {
			return "", err
		}
	}
} 

func existsAndIsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, fmt.Errorf("Directory %q does not exist. %w", path, err)
	} else if ! info.IsDir() {
		return false, fmt.Errorf("%q is not a directory.", path)
	} else {
		return true, nil
	}
}
	


// Asks the user for a Y/N answer; return true if Y or y, or false.
func answer() bool {
	fmt.Print(" Y/N: ")
	yn := false
	var answer string
	_, err := fmt.Scanln(&answer)
	if (err != nil) {
		fmt.Println("Error reading input. Assuming \"N\" (no/false).")
	} else if strings.EqualFold(string(answer[0]), "y") {
		yn = true
	}
	return yn
}

func flatten(dir string, prefix string, outdir string, partSep string) error {
	flatfunc := func(path string, info fs.DirEntry, err error) error {
		if (err != nil) {
			fmt.Printf("Error %v at path %q.", err, path)
		}
		if ! info.IsDir() {
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				fmt.Println("Cannot determine relative path.")
			} else {
				parts := splitPath(relPath) 
				newName := strings.Join(parts, partSep)
				oldPath := path // To do: This deems incorrect
				newPath := filepath.Join(outdir, newName)
				// b := copyFile(oldPath, newPath)
				// fmt.Printf("Copied %q to %q (%d bytes)\n", oldPath, newPath, b)
				copyFile(oldPath, newPath)
			}
		}
		return err
	}
	err2 := filepath.WalkDir(dir, flatfunc)
	return err2
}

func splitPath(path string) []string {
	sep := filepath.Join("a","b")[1]
	return strings.Split(filepath.Clean(path), string(sep))
}

func reverseStrings(original []string) []string {
	l := len(original)
	reversed := make([]string, l)
	for i:=0; i<l; i++ {
		reversed[i] = original[l-1-i]
	}
	return reversed
}

func copyFile(from, to string) int64 {
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
	return(r)
}


func usageText() string {
	return `
PURPOSE
Recursively copy and rename files from a "top directory" to one single
directory, prepending the new filenames with their origin path,
starting from the top directory. 
Files are read from the top directory (first argument) and copied
under their new names into the target directory. If no target directory
is specified (i.e. no second command line argument is given) then system
proposes alternatives and, if none is accepted, ends the program without
taking action.

CALL
   flatten [-sep <filename-separator>] <top-directory> [<output-path>]

EXAMPLE
Imagine the following directory with sub-directories and files
and note there are also empty folders included:

  Documents
      photos
       |--cats
       |   |--Molly
       |   |   |--seasick.jpg
       |   |   |--eating.jpg
       |   |   \--Molly-sleeps-on-balcony.jpg
       |   \--Jack
       |       |--sofanap.jpg
       |       |--riverboat.jpg
       |       \--chasing_mouse.tiff
       \--dogs
           \--Lassie
               |--portrait.jpg
               |--Lassie20250115.jpg
               \--Lassie-Drawing.png

Call: flatten Documents/photos 

This copies the contents of all files to the output directory, which is
by default the top directory (here Documents/photos) under new names,
which incorporate the previous path. The files would be named:
photos_cats_Molly_seasick.jpg
photos_cats_Molly_eating.jpg
photos_cats_Molly_Molly-sleeps-on-balcony.jpg
photos_cats_Jack_sofanap.jpg
...
photos_dogs_Lassie_portrait.jpg
photos_dogs_Lassie20250115.jpg
photos_dogs_Lassie-Drawing.png

PARAMETERS

-sep filename-separator
   The characters put between the previous subdirectory names.
   By default this is _, but could one or more character of any kind,
   but it is recommended not to seek for trouble with system
   conventions: Some characters like '*', '\', '$', '.', '~' and others
   have special meanings on popular computer systems and should
   therefore not be used in filenames.
   avoided.
   Recommended: "_" (default), "-", "--", "__".
   You could use the space character, though there are good reasons to
   not use spaces in filenames. In general, it is good practice to
   use no other characters in filenames than a-z, A-Z, 0-9, '-', '_',
   and the dot ('.') only between the actual name and the filename
   suffix; even if operating systems normally allow much more..
   With this flag set, you will see additional messages while the
   programme is running, e.g. which files get renamed etc. 

Testing:
Create and populate a directory ./testfiles 
(Launching creattestfiles.pl will do that for you.)
Then call the following to flatten the directory structure under
./testfiles into the directory ./testresults :
flatten ./testfiles ./testresults
`

}


package main

import(
	"fmt"
	"flag" 
	"log"
	"os"
	"path/filepath"
	"io"
	"strings"
	"errors"
)

var verbose bool   // Extra messages on actions taken

func main() {
	flag.Usage = func() {
		fmt.Println(usageText()) // No extra flag.PrintDefaults()
	}
	pathSepPtr := flag.String("sep", "_", "Path separator in output filenames")
	verbosePtr := flag.Bool("v", false, "Verbose messages during execution")
	silentPtr  := flag.Bool("s", false, "Silent mode - all output supporessed")
	flag.Parse()
	pathSep := *pathSepPtr
	verbose := *verbosePtr
	silent  := *silentPtr

	if verbose { fmt.Println("Verbose mode") }
	topdir, err := getTopDir(flag.Args(), silent)
	if err != nil {
		if ! silent { fmt.Println("Stopping here") }
		os.Exit(0)
	} else if ! silent {
		fmt.Printf("Top directory for flattening: %s\n", topdir)
	}

	outdir, err := getOutDir(topdir, flag.Args(), silent)
	if err == nil {
		if (verbose) {
			fmt.Println("Equivalent start prefix:     ", pathToPrefix(topdir))
			fmt.Println("Output directory:            ", outdir)
			fmt.Println("Path separator in new files: ", pathSep)
			fmt.Print(  "Is it ok to go ahead?  ")
			if (! answer()) {
				fmt.Println("=> Exiting")
				os.Exit(0)
			}
		}
		flatten(topdir, pathToPrefix(topdir), outdir, pathSep)
		if ! silent { fmt.Println("Done") }
	} else {
		fmt.Println("Ok, then we stop here.", err)
		os.Exit(0)
	}
}

func getOutDir(topdir string, args []string, silent bool) (string, error)  {
	var r string
	if len(args) > 1 {
		// User specified an output directory (argument 2).
		r = args[1]
	} else if silent {
		return "", errors.New("No output directory specified.")
	} else {
		// User did not specify an output directory.
		// As we are not in silent mode:
		//    - Ask to use top directory (wherever it comes from)
		//    - If not accepted, return an "error" message.
		fmt.Println("No output directory specified.")
		fmt.Printf("Should the output go to %s?", topdir)
		if answer() {
			r = topdir
		} else {
			return "", errors.New("Cannot write to nirvana.")
		}
	}
	return r, nil
} 

func getTopDir(args []string, silent bool) (string, error) {
	var r string
	if len(args) < 1 {
		if silent {
			return "", fmt.Errorf("Called without a top directory.")
		} else {
			fmt.Println("No top directory specified when calling.")
			fmt.Println("Picking current working directory.")
			r, err := os.Getwd()
			if err != nil {
				fmt.Println("Could not determine working directory.")
				return "", err
			} else {
				fmt.Printf("Working directory is: %q.\n", r)
				fmt.Print("Do you want to flatten that?")
				if answer() {
					return r, nil
				} else {
					return "", fmt.Errorf("Ok, no. So we flatten nothing.")
				}
			}
		}
	} else {
		r = args[0]
		_, err := os.Stat(r)  // Don't need the file information
		if os.IsNotExist(err) {
			if ! silent {
				fmt.Printf("You want to flatten %q but that seems not to exist.\n", r)
			}
			return "", fmt.Errorf("%q seems not to exist. %s", r, err)
		} else {
			return r, nil
		}
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

func flatten(dir string, prefix string, outdir string, sep string) {
	if verbose {
		fmt.Printf("Directory %s\n  => Prefix %s\n", dir, prefix)
	}
	entries, err := os.ReadDir(dir)
	if (err != nil) {
		fmt.Printf("Error reading directory %q.\n", dir)
		return
	}
	for _, fi := range entries {
		if fi.IsDir() {
			newpre := prefix + sep + fi.Name()
			from := filepath.Join(dir, fi.Name()) 
			flatten(from, newpre, outdir, sep)
		} else {
			// new file name with directory information included
			newfn := prefix + sep + fi.Name()
			// full path of the new (renamed) file:
			oldpath := filepath.Join(dir, fi.Name())
			newpath := filepath.Join(outdir, newfn)
			// fmt.Println("prefix: ", prefix)
			// fmt.Println("newpath: ", newpath)
			if verbose {
				fmt.Printf("  Copying file %s\n", fi.Name())
				fmt.Printf("    to %s\n", newpath)
			}
			// fmt.Println("copy", oldpath, "to", newpath)
			b := filecopy(oldpath, newpath)
			if verbose { fmt.Printf("    %d bytes written\n", b) }
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
       |   |   \--sleeping.jpg
       |   \--Jack
       |       |--sofanap.jpg
       |       |--riverboat.jpg
       |       \--chasing_mouse.tiff
       \--dogs
           |--Lassie
           |   |--Movie_Screenshots
           |   |   |--race2024.jpg
           |   |   \--trailerstart2024.jpg
           |   |--Images_2025
           |   \--Images_2026
           \--Sugar
               |--Holidays2020
                   |--Canada
                   |--Luxembourg
                       |--City
                           |--On_Palais_lawn.jpg
                           \--Jumping_from_bridge.jpg

Call: flatten Documents/photos 

This copies the contents of all files to the output directory, which is
by default the top directory (here Documents/photos) under new names,
which incorporate the previous path. The files would be named:
photos_cats_Molly_seasick.jpg
photos_cats_Molly_eating.jpg
photos_cats_Molly_sleeping.jpg
photos_cats_Jack_logo.png
...
photos_dogs_Sugar_Holidays2020_Luxembourg_City_On_Palais_lawn.jpg
photos_dogs_Sugar_Holidays2020_Luxembourg_City_Jumping_from_bridge.jpg

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
(tool creattestfiles.pl will do that for you).
Call: flatten ./testfiles ./testresults
`

}


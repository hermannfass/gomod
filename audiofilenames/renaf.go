package main

import(
	"log"
	"fmt"
	"os"
	"bufio"
	"strings"
	"path/filepath"
	"flag"
)

func main() {

	// Set the default Track Configuration Path (for user)
	// ---------------------------------------------------
	defConfPath, dcperr := checkDefConfPath()
	if dcperr != nil {
		log.Panic(dcperr)
	}
	confPathFlag := flag.String("c", defConfPath,
	                            "Path of Track Configuration File")
	fFlag := flag.Bool("f", false, "Force renaming without confirmation")
	flag.Usage = func() {
		fmt.Println(usageText(), "\nParameters:")
		flag.PrintDefaults()
		fmt.Println("")
	}
	flag.Parse()

	// Set the track path from argument 0
	// ----------------------------------
	var trackPath string
	if len(flag.Args()) == 0 {
		fmt.Println("You did not specify a target directory as parameter.")
		fmt.Println("Call `" + os.Args[0] + " -h` for usage instructions.")
		os.Exit(1)
	} else if _, tperr := os.Stat(flag.Args()[0]); os.IsNotExist(tperr) {
		fmt.Printf("The path %s does not exist.\n", flag.Args()[0])
		os.Exit(2)
	} 
	trackPath = flag.Args()[0]

	// Set the Track Configuration File (-c flag) and read it
	// ------------------------------------------------------
	confPath := *confPathFlag
	if _, cferr := os.Stat(confPath); os.IsNotExist(cferr) {
		fmt.Println("Track config file", confPath, "does not exist.")
		os.Exit(1)
	}
	ids, ierr := assignIDs(confPath) // read track conf into map
	if ierr != nil {
		fmt.Println("Could not evaluate the Track Configuratin File.")
		os.Exit(1)
	}

	// Rename files
	// ------------
	newfns, derr  := newFilenames(trackPath, ids)
	if (derr != nil) {
		log.Fatal("Problem reading files at %s.\n", derr)
	}
	var proceed bool
	if *fFlag {
		proceed = true
	} else {
		blessed, berr := blessing(trackPath, newfns)
		if berr != nil {
			log.Fatal("Problem getting blessing to proceed.\n", berr)
		}
		if blessed {
			proceed = true
			fmt.Println("Renaming files")
		} else {
			fmt.Println("Aborting.")
		}
	} 
	if proceed {
		for ofn, nfn  := range newfns {
			ofp := filepath.Join(trackPath, ofn)
			nfp := filepath.Join(trackPath, nfn)
			rerr := os.Rename(ofp, nfp)
			if rerr != nil {
				fmt.Printf("Could not rename %s to %s!\n", ofp, nfp)
				fmt.Println("Error:", rerr)
				continue
			}
			fmt.Printf("Renamed %s to %s.\n", ofn, nfn) 
		}
	}
}

// Confirm that we want to proceed
func blessing(tp string, newfns map[string]string) (bool, error) {
	fmt.Printf("%d files to rename in \"%s\".\n", len(newfns), tp)
	if len(newfns) == 0 {
		return false, nil
	}
	for old, new := range(newfns) {
		fmt.Printf("%20s ==> %-20s\n", old, new) 
	}
	fmt.Print("Do you want to rename those? [Y/n] ")
	r := bufio.NewReader(os.Stdin)
	c, cerr := r.ReadString('\n')
	if cerr != nil {
		return false, cerr
	}
	c = strings.TrimSpace(c)
	if c == "" || c == "Y" || c == "y" {
		return true, nil
	} else {
		return false, nil
	} 
}


// Return a map with names of all files to be renamed (key) and their 
// new filenames (value). Returns an error if it occurs.
func newFilenames(path string, ids map[string]string) (map[string]string, error) {
	trackfiles, derr := os.ReadDir(path)
	if derr != nil {
		return nil, derr
	}
	newfns := map[string]string{}  // New filenames; result of this function
	for _, trkf := range trackfiles {
		ofn := trkf.Name()  // Old filename
		ext := filepath.Ext(ofn)
		// Skip directories:
		if trkf.IsDir() {
			fmt.Printf("Skipping %s: It is a directory.\n", ofn) 
			continue
		}
		// See if the filename contains an ID and skip if it does not:
		nlabel := ""
		for instr, id := range ids {
			if strings.Contains(ofn, id) {
				nlabel = instr
			}
		}
		if nlabel == "" {
			fmt.Printf("Skipping %s: Contains no ID element\n", ofn)
			continue
		}
		// The file should get renamed. => Add it to the newNames map:
		nfn := nlabel + ext
		newfns[ofn] = nfn
	}
	return newfns, nil
}

// Read labels from config file into instr map
// func numToInstr(path string) (map[string]string, error) {
func assignIDs(path string) (map[string]string, error) {
	cf, ferr := os.Open(path)
	if (ferr != nil) {
		fmt.Printf("Could not open track configuration file %s.", path)
		return nil, ferr
	}
	defer cf.Close()
	scanner := bufio.NewScanner(cf)
	// No need to call scanner.Split(<func>) as ScanLines is default.
	ids := map[string]string{}  // Return value of this function
	for scanner.Scan() {
		cols := strings.Fields(scanner.Text())
		if len(cols) > 1 {
			ids[cols[1]] = cols[0]
		} else {
			// fmt.Printf("Skipping line: %s\n", scanner.Text())
		}
	}
	return ids, nil
}

// Determine default configuration path
func checkDefConfPath() (string, error) {
	if hd, uerr := os.UserHomeDir(); uerr != nil {
		fmt.Println("Looks like you have no home directory.")
		fmt.Println("Closing due to error:", uerr.Error())
		return "", uerr
	} else {
		return filepath.Join(hd, ".renaf", "tracks.conf"), nil
	}
}


// Usage Text
func usageText() string {
	return `
┌────────────────────────────────────────┐
│  RENAF - Rename Audio Recording Files  │
└────────────────────────────────────────┘

PURPOSE
Rename files in a directory according to a configuration table.
The configuration table is a list of identifiers (expected to be
present in a current filename) and their corresponding labels (that
will be used in the new filename). Filename extensions are preserved.
The original use case for this is naming numbered audio files from
Digital Audio Workstation. Filenames like "TRACK01.WAV", "TRACK02.WAV"
etc. get translated to person or instrument names that tell something
about the content, e.g. "Voice_Janis.WAV" or "Bassdrum.WAV" etc.

CALL
  renaf [-f] [-c path/to/track.conf] dir/with/files/to/rename",

TRACK CONFIGURATION FILE
The configuration file lists for each audio channel (track) a directory
unique identifier (in below example two-digit track numbers) and a more
descriptive file basename for the new filename. Identifier and new
basename are separated by whitespace.
Content of a typical RENAF configuration file:
  01  Voice_Janis
  02  Voice_Peter
  03  Guitar_James
  05  Guitar_Sam
  06  Bass
  07  Bassdrum
  08  Snare
	...
With this configuration, a file like "TRACK01.WAV" (note the "01" in
that filename) would get renamed into "Voice_Janis.WAV".

FILE LOCALIZATIONi
Ideally you place the configuration file(s) in the subdirectory ".renaf"
under your home directory and name "tracks.conf".
You can specify the path to the configuration file with the -c option.
For multiple projects, it is recommended to write separate configuration
files and copy the applicable one to "tracks.conf" (default) or to
overwrite that default path with the -c option.
`
}

package main

import(
	"log"
	"fmt"
	"os"
	"bufio"
	"strings"
	"regexp"
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
	instr, ierr := numToInstr(confPath) // read track conf into map
	if ierr != nil {
		fmt.Println("Could not evaluate the Track Configuratin File.")
		os.Exit(1)
	}

	// Rename files
	// ------------
	oldfn, newfn, derr  := fnForRen(trackPath, instr)
	if (derr != nil) {
		log.Fatal("Problem reading files at %s.\n", derr)
	}
	var proceed bool
	if *fFlag {
		proceed = true
	} else {
		blessed, berr := blessing(trackPath, oldfn, newfn)
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
		for i, ofn := range oldfn {
			nfn := newfn[i]
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
func blessing(tp string, oldfn []string, newfn []string) (bool, error) {
	fmt.Printf("%d files to rename in \"%s\".\n", len(oldfn), tp)
	if len(oldfn) == 0 {
		return false, nil
	}
	for i, f := range(oldfn) {
		fmt.Printf("%20s ==> %-20s\n", f, newfn[i])
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


// Return a hash with names of all fiels to be renamed
// and a hash with new filenames for those files, calculated
// from the track configuration in the instr map.
func fnForRen(path string, instr map[string]string) ([]string, []string, error) {
	trackfiles, derr := os.ReadDir(path)
	if derr != nil {
		return nil, nil, derr
	}
	re := regexp.MustCompile(`\d+`)
	var oldNames []string
	var newNames []string
	for _, trkf := range trackfiles {
		ofn := trkf.Name()  // Old filename
		// Skip directories:
		if trkf.IsDir() {
			fmt.Printf("Skipping %s: It is a directory.\n", ofn) 
			continue
		}
		// Skip files with no number in the filename:
		num := re.FindString(ofn)
		if num == "" {
			fmt.Printf("Skipping %s: No number in filename.\n", ofn) 
			continue
		}
		// Skip files with no label defined in track configuration:
		label, ok := instr[num]
		if !ok {
			fmt.Printf("Skipping %s: No label for track %s.\n", ofn, num)
			continue
		}
		// The file should get renamed.
		ext := filepath.Ext(ofn)  // File extension (usually ".WAV")
		nfn := label + ext
		oldNames = append(oldNames, ofn)
		newNames = append(newNames, nfn)
	}
	return oldNames, newNames, nil
}

// Read labels from config file into instr map
func numToInstr(path string) (map[string]string, error) {
	cf, ferr := os.Open(path)
	if (ferr != nil) {
		fmt.Printf("Could not open track configuration file %s.", path)
		return nil, ferr
	}
	defer cf.Close()
	scanner := bufio.NewScanner(cf)
	// No need to call scanner.Split(<func>) as ScanLines is default.
	instr := map[string]string{}
	for scanner.Scan() {
		cols := strings.Fields(scanner.Text())
		if len(cols) > 1 {
			// fmt.Printf("Fields: %q\n", cols)
			instr[cols[0]] = cols[1]
		} else {
			// fmt.Printf("Skipping line: %s\n", scanner.Text())
		}
	}
	return instr, nil
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
RENAF - Rename Audio Files

PURPOSE
Rename files in a directory according to a translation table specified
in a configuration file. The original use case for this is related to
naming numbered audio files from a Digital Audio Workstation. Filenames
there might look like "TRACK01.WAV", "TRACK02.WAV" etc. and the audio
engineer might want to see filenames like "Voice_Janis.WAV",
"Bassdrum.WAV" etc.

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

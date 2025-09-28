/*
Package songbook supports the merge of individual PDF files
into a PDF songbook, sorted by a playlist (text file) or by
alphabet. */
package songbook


import(
	"log"
	"fmt"
	"strings"
	"os"
	"regexp"
	"bufio"
	"path/filepath"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// essenceRE is a regular expression that describes the characters
// that will be taken out before matching a playlist entry (a title)
// against a PDF filename to decide if that PDF file is the sheet
// music for this title.
var essenceRE = regexp.MustCompile(`\W`)

// SongbookByList is core function 1/2:
// It compiles a songbook with sheet music sorted by a playlist
// file. Parameters are the FQFN of the playlist file, the path
// to the directory with PDF files, and the FQFN of the output file.
// If applicable, it returns a slice of warnings or other messages.
func SongbookByList(listPath, pdPath, outPath string) []string {
	var messages []string  // Do we need it? Nothing to add?
	var titles []string
	var allPdNames []string = GetAllPdNames(pdPath)
	var pdNames []string // Names of PDF files to include
	titles = ReadPlaylist(listPath)
	for _, t := range titles {
		pdNamesToAdd := PdNamesForTitle(t, allPdNames)
		if len(pdNamesToAdd) > 0 {
			pdNames = append(pdNames, pdNamesToAdd...)
		} else {
			messages = append(messages, "WARNING: No PDF file for " + t)
		}
	}
	MergePdfFiles(pdPath, pdNames, outPath)
	return messages
}

// SongbookByAbc is core function 2/2:
// It compiles an alphabetic songbook based on a PDF path and a
// playlist file path. If applicable, it returns a slice of warnings
// or other messages.
func SongbookByAbc(pdPath, outPath string) []string { 
	var messages []string
	var allPdNames []string = GetAllPdNames(pdPath)
	var pdNames []string  // Names of all files to be added
	for _, fn := range allPdNames {
		if okForAbcList(fn) {
			pdNames = append(pdNames, fn)
			fmt.Printf("Adding PDF file:   %s\n", fn)
		} else {
			fmt.Printf("Skipping PDF file: %s\n", fn)
		}
	}
	MergePdfFiles(pdPath, pdNames, outPath)
	return messages
}

// okForAbcList checks a filename and decides if the file should be
// considered when building an alphabetic songbook. This allows to
// name PDF title pages (like just indicating "Pause" or "Encores")
// with a prefix and thus keep them out of alphabetic songbooks.
func okForAbcList(fn string) bool {
	var excludedPrefixes = []string{"zzz", "ZZZ"}
	var ok bool = true
	for _, p := range excludedPrefixes {
		if strings.HasPrefix(fn, p) {
			ok = false
			break	
		}
	}
	return ok
}

// MergePdfFiles merges the files listed with their filenames in
// pdNames, located in the folder pdPath, into one PDF file that
// will be available at outPath.
func MergePdfFiles(pdPath string, pdNames []string, outPath string) {
	var pdps []string // PD Paths to combine
	for _, n := range pdNames {
		p := filepath.Join(pdPath, n)
		pdps = append(pdps, p)
	}
	if err := api.MergeCreateFile(pdps, outPath, false, nil); err != nil {
		fmt.Println("ERROR!")
		log.Fatal(err)
	}
}

// PdNamesForTitle returns a slice with one or more names of PDF
// filenames that match the title (song) in question.
func PdNamesForTitle(title string, fns []string) []string {
	var matches []string
	for _, fn := range fns {
		if fileMatch(fn, title) {
			matches = append(matches, fn)
		}
	}
	return matches
}

// ParseListName extracts project name and context from a
// playlist name and returns those two elements.
func ParseListName(list string) (string, string) {
	// Called only once, thus re not to be globally available.
	// Value is static, thus not using re, err = regexp.Compile().
	re := regexp.MustCompile(`(\w+)-(\w+)`)
	m := re.FindStringSubmatch(list)
	fmt.Printf("%v\n", m)
	return m[1], m[2] // m[0] is the whole match
}

// ReadPlaylist opens the text file (playlist) at the given
// path and returns a list (slice) of all lines. Empty lines
// and lines with only whitespace are ignored.
func ReadPlaylist(path string) []string {
	fh, err := os.Open(path)
	if (err != nil) {
		log.Fatal(err)
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	var entries []string
	for scanner.Scan() {
		line := scanner.Text();
		if strings.TrimSpace(line) == "" {
			continue
		}
		entries = append(entries, line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return entries
}

// GetAllPdNames takes a folder path and returns a list (slice) with
// the names of all the PDF files in this folder. PDF files are
// detected by the ".pdf" filename suffix.
func GetAllPdNames(path string) []string {
	var fns []string // List (Slice) of filenames to return
	// Regexp: Start with letter
	re := regexp.MustCompile(`(?i)\A[A-Za-z0-9][-\w]*\.pdf`)
	des, err := os.ReadDir(path) // DirectoryEntrys
	if (err != nil) {
		log.Fatal(err)
	}
	for _, de := range des { 
		fn := de.Name()
		if re.MatchString(fn) {
			fns = append(fns, fn)
		} else { 
			fmt.Printf("Skipping %s\n", fn) // debug
		}
	}
	return fns
}
	
// fileMatch reports whether a (PDF) filename contains to a certain
// extent a string (song title). Before this check, the two strings
// are prepared by the essence function to make the check case
// insensitive and to ignore special characters etc.
func fileMatch(filename, title string) bool {
	return strings.Contains(essence(filename), essence(title))
}

// essence extracts the meaningful parts of a string for
// a fuzzy match. It is used in the fileMatch function.
func essence(s string) string {
	n := strings.ToLower(essenceRE.ReplaceAllString(s, ""))
	return n
}



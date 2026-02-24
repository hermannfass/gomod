package main

import(
	"os"
	"path/filepath"
	"regexp"
	"fmt"
	"flag"
	"github.com/hermannfass/gomod/songbook"
)

// In this code, "pd" stands for portable document (PDF file).

// regexp to extract the "essence" from a song title
// and search it in the "essence" of a PDF filename.
// (Global, but limited to this package.)
var essenceRE = regexp.MustCompile(`\W`)

func main() {
	flag.Usage = func() {
		printUsageText()
		flag.PrintDefaults()
	}

	// Messages to be printed at the end of the execution;
	// mainly to list titles that have no matching PDF file.
	home := os.Getenv("HOME")
	// Default base path, overwritten with bpFlag:
	defaultBasePath := filepath.Join(home, "sheetmusic")
	bpFlag := flag.String("bp", defaultBasePath, "Base Path")
	listDirFlag := flag.String("lp", "playlists",
                  "Playlist Directory (relative to Base Path)")
	genPdDirFlag := flag.String("gendir", "Original",
	                "Name of directory with generic PDF files")
	flag.Parse()

	bp := *bpFlag
	listDir := filepath.Join(bp, *listDirFlag)
	fmt.Printf("Base path: %s  Playlist dir: %s\n", bp, listDir)

	listName := flag.Arg(0)
	project, context := songbook.ParseListName(listName)

	// Folder with the individual PDF files:
	pdPath := filepath.Join(bp, project)

	// Folder with generic PDF files:
	genPdPath := filepath.Join(bp, *genPdDirFlag)

	// Where to write the resulting songbook to:
	outName := fmt.Sprintf("%s-%s.pdf", project, context)
	outPath := filepath.Join(bp, outName)

	var messages []string

	if (context == "abc") {
		fmt.Printf("Collecting all PDF files from: %s\n", pdPath)
		fmt.Println("Compiling files sorted by alphabet.")
		fmt.Printf("Writing new songbook to: %s\n", outPath)
		messages = songbook.SongbookByAbc(pdPath, outPath) 
	} else {
		// to do: If arg looks like a path, take it as listPath?
		listPath := filepath.Join(listDir, listName)
		fmt.Printf("Reading sequence of repertoire from: %s\n", listPath)
		fmt.Printf("Collecting respective PDF files from: %s\n", pdPath)
		fmt.Printf("Writing new songbook to: %s\n", outPath)
		messages = songbook.SongbookByList(listPath, pdPath, genPdPath, outPath)
	}

	if len(messages) > 0 {
		fmt.Println("\nNOTE:")
		for _, m := range messages {
			fmt.Println(m)
		}
	}

}



func printUsageText() {
	fmt.Println(`
 /////////////
   songbook  
/////////////

PURPOSE

Combine PDF sheet music for a »Project«, that means for a band, an
orchestra, or a specific concert, into a »Songbook«, i.e. a single
PDF file.
For that, PDF files specific to the Project need to be kept in one
directory, called the »Project Folder«. A Project Folder may have
subdirectories, but they are ignored. 
The order in which the songs are added to a Songbook is defined
either by a »Playlist« or by alphabet, as described in the next
section. File organization details see further below and under
»File Naming and Localization«.

DEFINING CONTENT AND ORDER

The content (which pieces) and the order can be defined two ways:

 a) Songbook from a »Playlist«

   Playlists:
   A Playlist is a text file that lists one song title per line in
   the desired order. Its filename must consist of minimum two
   hyphen-separated parts and should end with ».txt« as suffix.
   <Project>-<Context>[-...].txt
   <Project> is the name of the Project (band, orchestra).
   <Context> represents the purpose for this list, like a specific
   concert, tour, or time period. This part must also not be empty.

   Playlist entries:
   Just list the songs that should be included in the Songbook one
   song per line. The entry for a song does not need to match the
   PDF filename exactly, but it must be included in the PDF filename.
   Matching Playlist entries against PDF filenames is case-insensitive
   and ignores spaces and non-alphanumeric character.
   For example, the PDF file »BeautifulNoise-NeilDiamond-guitar.pdf«
   will match a Playlist entry like »Beautiful Noise« or
   »Beautiful Noise (Neil Diamond)«. Note that the Playlist entry
   needs to be part of the filename, not vice versa.

   Avoid amiguity:
   When entering songs in a Playlist, make sure they point to only
   one PDF file; otherwise your Songbook might not include what you
   expect. You may have more versions of one song or separate PDF files
   for different instruments. Thus, the Playlist entry needs to be
   specific enough. If there are two files for Gershwin's »Summertime«,
   like »Summertime-BigBrother.pdf« and »Summertime-Holiday.pdf«, a
   Playlist entry like just »Summertime« is not clear, but for example
   »Summertime (Holiday)« will be clear.
   Also note that a Playlist entry like »Summertime (Billie Holiday)«
   will not match any of the two PDF files: The system removes spaces
   and non-alphanumeric characters from the Playlist entry and from all
   filenames, downcases the results and then looks for a stripped
   filename that includes the stripped Playlist entry,
   »summertimebillieholiday«. Neither »summertimeholday.pdf« nor
   »summertimebigbrother.pdf« includes this string, so there will be
   no result, i.e. no PDF file for this song will go to the Songbook.

   Cross-Project PDF files:
   If no file in the Project Folder seems to match the name of the
   song title, the application looks for a match in a folder with
   generic sheet music, i.e. original versions that can be used in
   multiple Projects. If the file is taken from there, this is
   indicated in the output during the run.

   Comments:
   In case you want to add comments to your Playlist (or make the
   system temporarily ignoring individual entries in the Playlist),
   just prepend the respective lines with a hash symbol (#).
   
   Example:
   Assuming one Playlist for project »CoolBand« is called
   »CoolBand-Concert20250913.txt«
   To create a songbook for this event, call the tool like this:
   songbook CoolBand-Concert20250913.txt
   This will create a PDF file called CoolBand-Concert20250913.pdf
   in the output folder (location see below under
   »File Naming and Localization«).

 b) Songbook by alphabet:

   If you just want all the PDF files in a Project Folder, instead
   of specifying a playlist, you just use the Project name followed
   by a hyphen and the three letters »abc«.
   Example: songbook CoolBand-abc
   This will combine all PDF files in the Project Folder »CoolBand«
   into one PDF file named »CoolBand-abc.pdf«.

FILE NAMING AND LOCALIZATION

   In general, it is good style (not only for this application) to
   use only letters, numbers, underscores, and hyphens in folder and
   file names; no spaces, umlauts, accents, or other special letters.

   Base Path:
   Parent directory for all PDF files and playlists. Its default
   value is the folder called »sheetmusic« in the user's home
   directory, so technically: ~/sheetmusic
   You can change the default value with the respective flag
   (see PARAMETERS below).
   
   Project Folders:
   Under the Base Path, each Project (band, orchestra) must have a
   folder that contains the sheet music for this project, i.e. an
   individual PDF file for each piece in the repertoire of this
   ensemble.

   Cross-Project Folder:
   Under the Base Path you may create a folder that contains sheet
   music which is (or might be) used in multiple projects. Usually
   these are standard or original versions without project specific
   modifications.
   The default location and name for this folder is:
   ~/sheetmusic/Original/

   Playlist Folder:
   Also under the Base Path should be one folder called "playlists",
   which contains Playlist Files. We assume you want to keep the
   playlists for all Projects in this one Playlist folder, though
   you can override it with the respective flag (see PARAMETERS below).

   Example for a typical file structure:

   Let's assume a person plays in a band called  »The Keltners« and
   is also booked for a one-time event »Rockstar Summit 2025«.
   Accepting the default values, this is how files should be
   organized. (Note: "~" stands for the user's home directory.)

      ~/sheetmusic (= Base Path)
      ~/sheetmusic/TheKeltners/               (= Project Folder 1)
      ~/sheetmusic/TheKeltners/Shalala.pdf
      ~/sheetmusic/TheKeltners/Ohyeah-v2-2021.pdf
      ~/sheetmusic/TheKeltners/WhateverYouWant.pdf
      ...
      ~/sheetmusic/RockstarSummit_2025        (= Project Folder 2)
      ~/sheetmusic/RockstarSummit_2025/Opening-arrangement.pdf
      ~/sheetmusic/RockstarSummit_2025/WeSaluteYou-Malcolm-guitar.pdf
      ~/sheetmusic/RockstarSummit_2025/TheBoxer-guitar.pdf
      ...
      ~/sheetmusic/Original                   (= Cross-Project Folder)
      ~/sheetmusic/Original/AutumnLeaves-LeadSheet.pdf
      ~/sheetmusic/Original/Toto-Africa-Drums.pdf
      ...
      ~/sheetmusic/playlists/TheKeltners-shortSet2025.txt
      ~/sheetmusic/playlists/TheKeltners-fullSet2025.txt
      ~/sheetmusic/playlists/RockstarSummit_2025.txt

   Call: songbook TheKeltners-shortSet2025.txt
      This will create a PDF with the songs as listed in the
      Playlist File named TheKeltners-shortSet2025.txt.
      If a Playlist Entry is not included in Project Folder 1,
      »~/sheetmusic/TheKeltners/«, the system will also look in the
      Cross-Project Folder, »~/sheetmusic/Original«.

   Call: songbook TheKeltners-abc
      This will create a PDF with all songs in Project Folder 1,
      listed by alphabet.

PARAMETERS
`)
}



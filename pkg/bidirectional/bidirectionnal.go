package bidirectional

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/tools/godoc/util"
	"golang.org/x/tools/godoc/vfs"

	"github.com/ariary/TrojanSourceFinder/pkg/config"
	"github.com/ariary/TrojanSourceFinder/pkg/excludelist"
	"github.com/ariary/TrojanSourceFinder/pkg/utils"
)

//Global dict for bidirectional characrters
var BidirectionalCharactersDict = map[rune]string{
	0x200E: "LRM",
	0x200F: "RLM",
	0x061C: "ALM",
	0x202A: "LRE",
	0x202B: "RLE",
	0x202C: "PDF",
	0x202D: "LRO",
	0x202E: "RLO",
	0x2066: "LRI",
	0x2067: "RLI",
	0x2068: "FSI",
	0x2069: "PDI",
}

var (
	prefix = "("
	suffix = ")"
)

// //Implicit Directional Formatting Characters
// var ImplicitDirectionalDict = map[rune]string{
// 	0x200E: "LRM",
// 	0x200F: "RLM",
// 	0x061C: "ALM",
// }

// //Explicit Directional Embedding and Override Formatting Characters
// var ExplicitDirectionalDict = map[rune]string{
// 	0x202A: "LRE",
// 	0x202B: "RLE",
// 	0x202C: "PDF",
// 	0x202D: "LRO",
// 	0x202E: "RLO",
// }

// //Explicit Directional Isolate Formatting Characters
// var ExplicitDirectionalIsolateDict = map[rune]string{
// 	0x2066: "LRI",
// 	0x2067: "RLI",
// 	0x2068: "FSI",
// 	0x2069: "PDI",
// }

// IsBidirectionalAlgorithm reports whether the rune is a bidirectionnal character as defined
// by Unicode's bidirectional algorithm property;
// This is:
// - Implicit Directional Formatting Characters: U+200E (LRM), U+200F (RLM), U+061C (ALM)
// - Explicit Directional Embedding and Override Formatting Characters: U+202A (LRE), U+202B (RLE), U+202D (LRO), U+202E (RLO), U+202C (PDF)
// - Explicit Directional Isolate Formatting Characters: U+2066 (LRI), U+2067 (RLI), U+2068 (FSI), U+2069 (PDI).
func IsBidirectionalAlgorithm(r rune) bool {
	if _, isIn := BidirectionalCharactersDict[r]; isIn {
		return true
	}

	return false
}

// ContainBidirectionnal reports whether the byte contains  bidirectionnal character as defined
// by Unicode's bidirectional algorithm property; this could lead to Trojan Source vulnerability
func ContainBidirectionnal(str string) bool {
	for _, c := range str {
		if IsBidirectionalAlgorithm(c) {
			return true
		}
	}
	return false
}

// Return the evil line with Bidirectional character replace
func getEvilLine(str string, color bool) (exorcisedStr string) {
	for _, c := range str {
		if s, isIn := BidirectionalCharactersDict[c]; isIn {
			//add bad unicode character with its representation
			s = prefix + s + suffix
			if color {
				exorcisedStr += utils.Bold(utils.RedForeground(s))
			} else {
				exorcisedStr += s
			}
		} else {
			exorcisedStr += string(c)
		}
	}
	return exorcisedStr
}

// Scan file or folder to detect potential Trojan Source vulnerability within.
// This function exit with status code 1 if trojan source has been detected, 0 otherwise
func Scan(path string, cfg *config.Config) {
	utils.InitLoggers()
	// Recursive (directory) or normal scan?
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	excludedPaths, err := excludelist.GetExcludelist(cfg.ExcludelistFilename)
	if err != nil {
		log.Fatal(err)
	}

	// Skip the given path if it is contained in the exclude list.
	if _, ignoreEntry := (*excludedPaths)[filepath.Clean(path)]; ignoreEntry {
		os.Exit(0)
	}

	var detected int
	if fileInfo.IsDir() {
		detected = scanDirectory(path, cfg, excludedPaths)
	} else {
		detected = scanFile(path, cfg)
	}

	os.Exit(detected)
}

// Scan a file to detect the presence of potential Trojan Source
// return 0 if no trojan source has been detected within file
func scanFile(filename string, cfg *config.Config) int {
	// test if human readable text
	if cfg.OnlyText {
		fs := vfs.OS(".")
		if !util.IsTextFile(fs, filename) { //Not a "human readable" file so probably not surce code
			if cfg.Verbose {
				result := "not scanned (not a text file)"
				if cfg.Color {
					result = utils.Italic(utils.Yellow(result))
				}
				utils.ErrorLogger.Println("check", filename, "...", result)
			}
			return 1
		}
	}

	/*SCAN*/
	detected := false
	line := 1
	vulns := make(map[int]string)

	// Reade file
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	//Increase max line length that could be read by scanner (<=1MB)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() { // read file line by line
		lineText := scanner.Text()
		lineDetector := ContainBidirectionnal(lineText)

		if lineDetector {
			detected = true
			vulns[line] = lineText
		}
		line++
	}
	if err := scanner.Err(); err != nil {
		if err == bufio.ErrTooLong {
			utils.ErrorLogger.Println(filename, ": too long lines (<1M), could not be checked")
		} else {
			log.Fatal(err)
		}
	}

	/*REPORT*/
	var result string
	if detected {

		if cfg.Color {
			result = utils.Evil("not ok")
		} else {
			result = "not ok"
		}
		utils.ErrorLogger.Println("check", filename, "...", result)
		if cfg.Verbose {
			for line, text := range vulns {
				msg := getEvilLine(text, cfg.Color)
				var lineS string
				lineS = strconv.Itoa(line)
				if cfg.Color {
					lineS = utils.Yellow(lineS)
				}
				utils.InfoLogger.Println(lineS, ": ", msg)
			}
		}
		return 1
	} else {
		if cfg.Verbose {
			if cfg.Color {
				result = utils.Green("ok")
			} else {
				result = "ok"
			}
			utils.InfoLogger.Println("check", filename, "...", result)
		}
	}
	return 0
}

// Scan recursively all the file of a repository (pathD) to detect the presence of
//potential Trojan Source.
// Browse the directory using filepath.Walk package => does not follow symbolic link
// and for very large directories Walk can be inefficient
// return 0 if no trojan source was detected
func scanDirectory(pathD string, cfg *config.Config, excludedPaths *map[string]struct{}) (result int) {
	err := filepath.WalkDir(pathD, func(path string, dirEntry os.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if dirEntry.IsDir() {
			// Skip the whole directory if it is contained in the exclude list.
			if _, ignoreEntry := (*excludedPaths)[path]; ignoreEntry {
				return filepath.SkipDir
			}
		} else {
			// Skip the file if it is contained in the exclude list.
			if _, ignoreEntry := (*excludedPaths)[path]; ignoreEntry {
				return nil
			}

			result += scanFile(path, cfg)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	// "normalize" result
	if result > 0 {
		result = 1
	}

	return result
}

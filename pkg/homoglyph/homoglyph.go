package homoglyph

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ariary/TrojanSourceFinder/pkg/config"
	"github.com/ariary/TrojanSourceFinder/pkg/excludelist"
	"github.com/ariary/TrojanSourceFinder/pkg/utils"
	confusable "github.com/skygeario/go-confusable-homoglyphs"
	"golang.org/x/tools/godoc/util"
	"golang.org/x/tools/godoc/vfs"
)

var (
	prefix = "["
	suffix = "]"
)

//Return All homoglyph in string
func getAllHomoglyph(str string) []string {
	var homoglyphes []string
	words := strings.Fields(str)
	for _, word := range words {
		if confusable.IsDangerous(word, []string{}) {
			homoglyphes = append(homoglyphes, word)
		}
	}
	return homoglyphes
}

// Return the evil line with Homoglyph in parenthesis. Homoglyph detection is done
//"word by word"
func getEvilLine(str string, color bool) (exorcisedStr string) {
	words := strings.Fields(str)
	for _, word := range words {
		if confusable.IsDangerous(word, []string{}) {
			word = prefix + word + suffix
			if color {
				word = utils.Bold(utils.RedForeground(word))
			}
		}
		exorcisedStr += word

	}
	return exorcisedStr
}

// Scan file or folder to detect potential homoglyph within.
// This function exit with status code 1 if homoglyph has been detected, 0 otherwise
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

// Scan a file to detect the presence of potential Homoglyph
// return 0 if no homoglyph has been detected within file
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
			return 0
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
		lineDetector := confusable.IsDangerous(lineText, []string{})

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
				utils.InfoLogger.Println(line, ": ", msg)
			}
		}
		/*SIBLING REPORT*/
		if cfg.Sibling != nil {
			for _, text := range vulns {
				for i := 0; i < len(*cfg.Sibling); i++ {
					path := (*cfg.Sibling)[i]
					getSiblings(path, text, cfg)
				}
			}
		}
		return 1
	} else {
		if cfg.Color {
			result = utils.Green("ok")
		} else {
			result = "ok"
		}
		utils.InfoLogger.Println("check", filename, "...", result)
	}
	return 0
}

// Scan recursively a repository to detect the presence of potential Homoglyph
// Browse the directory using filepath.Walk package => does not follow symbolic link
// and for very large directories Walk can be inefficient
func scanDirectory(filename string, cfg *config.Config, excludedPaths *map[string]struct{}) (result int) {
	err := filepath.WalkDir(filename, func(path string, dirEntry os.DirEntry, err error) error {
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

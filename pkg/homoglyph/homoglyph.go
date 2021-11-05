package homoglyph

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ariary/TrojanSourceFinder/pkg/utils"
	confusable "github.com/skygeario/go-confusable-homoglyphs"
)

var (
	prefix = "["
	suffix = "]"
)

// Return the evil line with Homoglyph in parenthesis. Homoglyph detection is done
//"word by word"

func getEvilLine(str string, color bool) (exorcisedStr string) {
	words := strings.Fields(str)
	for _, word := range words {
		if confusable.IsDangerous(word, []string{}) {
			word = prefix + word + suffix
			if color {
				exorcisedStr += utils.Bold(utils.RedForeground(word))
			}

		}
		exorcisedStr += word
	}
	return exorcisedStr
}

// Scan file or folder to detect potential homoglyph within.
func Scan(filename string, recursive bool, verbose bool, color bool) {
	utils.InitLoggers()

	if recursive {
		scanDirectory(filename, verbose, color)
	} else {
		scanFile(filename, verbose, color)
	}
}

// Scan a file to detect the presence of potential Homoglyphe
func scanFile(filename string, verbose bool, color bool) {
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

		if color {
			result = utils.Evil("not ok")
		} else {
			result = "not ok"
		}
		utils.ErrorLogger.Println("check", filename, "...", result)
		if verbose {
			for line, text := range vulns {
				msg := getEvilLine(text, color)
				utils.InfoLogger.Println(line, ": ", msg)
			}
		}
	} else {
		if color {
			result = utils.Green("ok")
		} else {
			result = "ok"
		}
		utils.InfoLogger.Println("check", filename, "...", result)
	}
}

// Scan recursively a repository to detect the presence of potential Homoglyph
// Browse the directory using filepath.Walk package => does not follow symbolic link
// and for very large directories Walk can be inefficient
func scanDirectory(filename string, verbose bool, color bool) {
	err := filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !info.IsDir() {
			scanFile(path, verbose, color)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

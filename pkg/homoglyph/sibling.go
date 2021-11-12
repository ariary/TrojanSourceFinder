package homoglyph

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ariary/TrojanSourceFinder/pkg/config"
	"github.com/ariary/TrojanSourceFinder/pkg/utils"
	"github.com/eskriett/confusables"
	mconfusables "github.com/mtibben/confusables"
)

// Return the sibling of a specific homoglyphe found in file. Homoglyph -> Skeleton. Search for word with same skeleton.
func getSiblings(path string, homoglyphLine string, cfg *config.Config) {
	if cfg.Verbose {
		utils.InfoLogger.Println("Search sibling homograph for line:", homoglyphLine)
	}
	// Recursive (directory) or normal scan?
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	if fileInfo.IsDir() {
		getSiblingsDirectory(path, homoglyphLine, cfg.Color, cfg.Verbose)
	} else {
		getSiblingsFile(path, homoglyphLine, cfg.Color, cfg.Verbose)
	}
}

// Return the sibling of a specific homoglyphe found in file. Homoglyph -> Skeleton. Search for word with same skeleton.
// To find sibling, check with 2 external libraires: mtibben/confusables & eskriett/confusables
// These libraries does not work very well
func getSiblingsFile(path string, homoglyphLine string, color bool, verbose bool) {
	homoglyphes := getAllHomoglyph(homoglyphLine)
	for _, homoglyph := range homoglyphes {
		//SCAN
		detected := false
		line := 1
		siblings := make(map[int][]string)
		// Reade file
		f, err := os.Open(path)
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
			words := strings.Fields(lineText)
			for _, word := range words {
				lineDetector := confusables.IsConfusable(homoglyph, word) || mconfusables.Confusable(homoglyph, word)
				if lineDetector {
					detected = true
					siblings[line] = append(siblings[line], word)
				}
			}
			line++
		}

		//REPORT
		if detected {
			if color {
				path = utils.Bold(utils.Magenta(path))
			}
			utils.ErrorLogger.Println("Find sibling in", path, ":")
			for line, sibling := range siblings {
				for i := 0; i < len(sibling); i++ {
					var lineS string
					lineS = strconv.Itoa(line)
					if color {
						lineS = utils.Yellow(lineS)
						sibling[i] = utils.Purple(sibling[i])
					}
					utils.ErrorLogger.Println("\t", lineS, ":", sibling[i])
				}
			}

		}
	}
}

// Return the sibling of a specific homoglyph found in directory
func getSiblingsDirectory(pathD string, homoglyphLine string, color bool, verbose bool) {
	err := filepath.Walk(pathD, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !info.IsDir() {
			getSiblingsFile(path, homoglyphLine, color, verbose)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

package bidirectional

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

//LOGGER
var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func initLoggers() {
	InfoLogger = log.New(os.Stdout, "", 0)
	ErrorLogger = log.New(os.Stderr, "", 0)
}

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
	// if _, isIn := ImplicitDirectionalDict[r]; isIn {
	// 	return true
	// }

	// //Explicit Directional Embedding and Override Formatting Characters
	// if _, isIn := ExplicitDirectionalDict[r]; isIn {
	// 	return true
	// }

	// //Explicit Directional Isolate and Override Formatting Characters
	// if _, isIn := ExplicitDirectionalIsolateDict[r]; isIn {
	// 	return true
	// }
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
func getEvilLine(str string) (exorcisedStr string) {
	for _, c := range str {
		if s, isIn := BidirectionalCharactersDict[c]; isIn {
			exorcisedStr += s
		} else {
			exorcisedStr += string(c)
		}
	}
	return exorcisedStr
}

// Scan file or folder to detect potential Trojan Source vulnerability within.
func Scan(filename string, recursive bool, verbose bool) {
	initLoggers()

	if recursive {
		scanDirectory(filename, verbose)
	} else {
		scanFile(filename, verbose)
	}
}

// Scan a file to detect the presence of potential Trojan Source
func scanFile(filename string, verbose bool) {
	/*SCAN*/
	detected := false
	line := 1
	//vulns := make(map[int]bidi.Ordering)
	vulns := make(map[int]string)

	// Reade file
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	//Increase max line length taht could be read by scanner (<=1MB)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() { // read file line by line
		//lineDetector, rOrd := ContainBidirectionnal(scanner.Bytes())
		lineText := scanner.Text()
		lineDetector := ContainBidirectionnal(lineText)
		if lineDetector {
			detected = true
			//vulns[line] = rOrd
			vulns[line] = lineText
		}
		line++
	}
	if err := scanner.Err(); err != nil {
		if err == bufio.ErrTooLong {
			ErrorLogger.Println(filename, ": too long lines (<1M), could not be checked")
		} else {
			log.Fatal(err)
		}
	}

	/*REPORT*/
	if detected {
		ErrorLogger.Println("check", filename, "... not ok")
		for line, text := range vulns {
			if verbose {
				msg := getEvilLine(text)
				InfoLogger.Println(line, ": ", msg)
			}
			// if exorcise {
			// 	msg := getExorcisedLine(ord)
			// 	InfoLogger.Println(line, ": ", msg)
			// }
		}
	} else {
		InfoLogger.Println("check", filename, "... ok")
	}
}

// Scan recursively a repository to detect the presence of potential Trojan Source
// Browse the directory using filepath.Walk package => does not follow symbolic link
// and for very large directories Walk can be inefficient
func scanDirectory(filename string, verbose bool) {
	err := filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !info.IsDir() {
			scanFile(path, verbose)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

// ContainBidirectionnal reports whether the byte contains  bidirectionnal character as defined
// by Unicode's bidirectional algorithm property; this could lead to Trojan Source vulnerability
// func ContainBidirectionnal(b []byte) (detected bool, ord bidi.Ordering) {
// 	var p bidi.Paragraph
// 	p.SetBytes(b)
// 	ord, err := p.Order()
// 	if err != nil {
// 		panic(err)
// 	}

// 	//Reconstruct Ordering.directions as it is a private fields
// 	rOrd := reflect.ValueOf(ord)
// 	rDirection := rOrd.FieldByName("directions")

// 	for i := 1; i < rDirection.Len(); i++ { // check if we have the same directions for all runes, In fact we could just test if len > 1
// 		if rDirection.Index(i) != rDirection.Index(0) {
// 			return true, ord
// 		}
// 	}
// 	return false, ord
// }

// Return the line contains within runes of bidi.Ordering object without dealing  w/
// reordering bidirectional characters
// func getEvilLine(ord bidi.Ordering) (msg string) {
// 	//Reconstruct Ordering.directions as it is a private fields
// 	rOrd := reflect.ValueOf(ord)
// 	rRunes := rOrd.FieldByName("runes")
// 	var rs []rune                       //will contain all concated runes
// 	for i := 0; i < rRunes.Len(); i++ { // check if we have the same directions for all runes, In fact we could just test if len > 1
// 		rRunesSlice := rRunes.Index(i)
// 		for j := 0; j < rRunesSlice.Len(); j++ { //does not find other way to reconstruc string
// 			r := rune(reflect.ValueOf(rRunesSlice.Index(j).Int()).Int())
// 			rs = append(rs, r)
// 		}
// 	}
// 	msg = string(rs)
// 	return msg
// }

// Return the line contains within runes of bidi.Ordering object and reorder
// bidirectional characters
// func getExorcisedLine(ord bidi.Ordering) (msg string) {
// 	// what is the default order ?
// 	//todo (from now our scope is Left-to-Right)
// 	direction := 0

// 	// indexes of wrong direction runes?
// 	var evilIndexes []int
// 	rOrd := reflect.ValueOf(ord)
// 	rDirection := rOrd.FieldByName("directions")

// 	for i := 0; i < rDirection.Len(); i++ { // check if the direction is the same as the default one
// 		if int(rDirection.Index(i).Int()) != direction {
// 			evilIndexes = append(evilIndexes, i)
// 		}
// 	}

// 	// reconstruct string, if index is in the wrong list reverseString
// 	rRunes := rOrd.FieldByName("runes")
// 	var rs []rune
// 	for i := 0; i < rRunes.Len(); i++ {
// 		rRunesSlice := rRunes.Index(i)
// 		var tmpS string
// 		var tmpR []rune
// 		for j := 0; j < rRunesSlice.Len(); j++ { //reconstruct string
// 			r := rune(reflect.ValueOf(rRunesSlice.Index(j).Int()).Int())
// 			tmpR = append(tmpR, r)
// 			rs = append(rs, r)
// 		}

// 		// determine if it is an evil
// 		evil := false
// 		for _, index := range evilIndexes {
// 			if index == i {
// 				evil = true
// 				break
// 			}
// 		}
// 		if evil {
// 			tmpS += bidi.ReverseString(string(tmpR))
// 		} else {
// 			tmpS += string(tmpR)
// 		}
// 		msg += tmpS
// 	}
// 	//msg = string(rs)
// 	return msg
// }

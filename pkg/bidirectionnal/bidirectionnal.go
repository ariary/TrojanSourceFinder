package bidirectionnal

import (
	"bufio"
	"log"
	"os"
	"reflect"

	"golang.org/x/text/unicode/bidi"
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

//Implicit Directional Formatting Characters
var ImplicitDirectionalDict = map[rune]string{
	0x200E: "LRM",
	0x200F: "RLM",
	0x061C: "ALM",
}

//Explicit Directional Embedding and Override Formatting Characters
var ExplicitDirectionalDict = map[rune]string{
	0x202A: "LRE",
	0x202B: "RLE",
	0x202C: "PDF",
	0x202D: "LRO",
	0x202E: "RLO",
}

//Explicit Directional Isolate Formatting Characters
var ExplicitDirectionalIsolateDict = map[rune]string{
	0x2066: "LRI",
	0x2067: "RLI",
	0x2068: "FSI",
	0x2069: "PDI",
}

// IsBidirectionalAlgorithm reports whether the rune is a bidirectionnal character as defined
// by Unicode's bidirectional algorithm property;
// This is:
// - Implicit Directional Formatting Characters: U+200E (LRM), U+200F (RLM), U+061C (ALM)
// - Explicit Directional Embedding and Override Formatting Characters: U+202A (LRE), U+202B (RLE), U+202D (LRO), U+202E (RLO), U+202C (PDF)
// - Explicit Directional Isolate Formatting Characters: U+2066 (LRI), U+2067 (RLI), U+2068 (FSI), U+2069 (PDI).
// func IsBidirectionalAlgorithm(r rune) bool {
// 	if code, isIn := ImplicitDirectionalDict[r]; isIn {
// 		fmt.Println("Detect Implicit directional formatting characters:", code)
// 		return true
// 	}

// 	//Explicit Directional Embedding and Override Formatting Characters
// 	if code, isIn := ExplicitDirectionalDict[r]; isIn {
// 		fmt.Println("Detect Explicit directional embedding and override formatting characters:", code)
// 		return true
// 	}

// 	//Explicit Directional Isolate and Override Formatting Characters
// 	if code, isIn := ExplicitDirectionalIsolateDict[r]; isIn {
// 		fmt.Println("Detect Explicit directional isolate and override formatting characters:", code)
// 		return true
// 	}
// 	return false
// }

// ContainBidirectionnal reports whether the byte contains  bidirectionnal character as defined
// by Unicode's bidirectional algorithm property; this could lead to Trojan Source vulnerability
func ContainBidirectionnal(b []byte) (detected bool, ord bidi.Ordering) {
	var p bidi.Paragraph
	p.SetBytes(b)
	ord, err := p.Order()
	if err != nil {
		panic(err)
	}

	//Reconstruct Ordering.directions as it is a private fields
	rOrd := reflect.ValueOf(ord)
	rDirection := rOrd.FieldByName("directions")

	for i := 1; i < rDirection.Len(); i++ { // check if we have the same directions for all runes, In fact we could just test if len > 1
		if rDirection.Index(i) != rDirection.Index(0) {
			return true, ord
		}
	}
	return false, ord
}

// Return the line contains within runes of bidi.Ordering object without dealing  w/
// reordering bidirectional characters
func getEvilLine(ord bidi.Ordering) (msg string) {
	//Reconstruct Ordering.directions as it is a private fields
	rOrd := reflect.ValueOf(ord)
	rRunes := rOrd.FieldByName("runes")
	var rs []rune                       //will contain all concated runes
	for i := 1; i < rRunes.Len(); i++ { // check if we have the same directions for all runes, In fact we could just test if len > 1
		rRunesSlice := rRunes.Index(i)
		for j := 0; j < rRunesSlice.Len(); j++ { //does not find other way to reconstruc string
			r := rune(reflect.ValueOf(rRunesSlice.Index(j).Int()).Int())
			rs = append(rs, r)
		}
	}
	msg = string(rs)
	return msg
}

// Return the line contains within runes of bidi.Ordering object and reorder
// bidirectional characters
func getExorcisedLine(ord bidi.Ordering) (msg string) {
	return msg
}

// Scan file or folder to detect potential Trojan Source vulnerability within.
// recursive: scan folder; exorcised: print out the vulnerability detected
func Scan(filename string, recursive bool, exorcise bool, verbose bool) {
	initLoggers()

	//SCAN
	detected := false
	line := 0
	vulns := make(map[int]bidi.Ordering)

	// Reade file
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// read file line by line (bidi.Paragraph: " If text contains a paragraph separator it
	// will only process the first paragraph")
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		lineDetector, rOrd := ContainBidirectionnal(scanner.Bytes())
		if lineDetector {
			detected = true
			vulns[line] = rOrd
		}
		line++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	//REPORT
	if detected {
		ErrorLogger.Println("check", filename, "... not ok")
		for line, ord := range vulns {
			if verbose {
				msg := getEvilLine(ord)
				InfoLogger.Println(line, ": ", msg)
			}
			if exorcise {
				msg := getExorcisedLine(ord)
				InfoLogger.Println(line, ": ", msg)
			}
		}
	} else {
		InfoLogger.Println("check", filename, "... ok")
	}

}

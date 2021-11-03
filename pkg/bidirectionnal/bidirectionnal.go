package bidirectionnal

import "fmt"

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
func IsBidirectionalAlgorithm(r rune) bool {
	if code, isIn := ImplicitDirectionalDict[r]; isIn {
		fmt.Println("Detect Implicit directional formatting characters:", code)
		return true
	}

	//Explicit Directional Embedding and Override Formatting Characters
	if code, isIn := ExplicitDirectionalDict[r]; isIn {
		fmt.Println("Detect Explicit directional embedding and override formatting characters:", code)
		return true
	}

	//Explicit Directional Isolate and Override Formatting Characters
	if code, isIn := ExplicitDirectionalIsolateDict[r]; isIn {
		fmt.Println("Detect Explicit directional isolate and override formatting characters:", code)
		return true
	}
	return false
}

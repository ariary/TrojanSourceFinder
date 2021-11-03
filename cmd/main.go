// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/ariary/TrojanSourceFinder/pkg/bidirectional"
)

func Example_is() {

	// constant with mixed type runes
	const mixed = "/*‮ } ⁦if (isAdmin)⁩ ⁦ begin admins only */"
	for _, c := range mixed {
		fmt.Printf("For %q:\n", c)
		if bidirectional.IsBidirectionalAlgorithm(c) {
			fmt.Println("\tis control rune")
		}
	}
}

func main() {
	Example_is()
}

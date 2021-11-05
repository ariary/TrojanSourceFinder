// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/ariary/TrojanSourceFinder/pkg/bidirectional"
	"github.com/spf13/cobra"
)

func main() {
	//CMD SCAN
	var recursive bool
	var verbose bool
	var color bool

	var rootCmd = &cobra.Command{
		Use:   "tsFinder [filename]",
		Short: "Detect Trojan Source Vulnerability in your code",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			name := args[0]

			bidirectional.Scan(name, recursive, verbose, color)
		},
	}

	//flag handling
	rootCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", false, "scan all the files in the specified folder")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "make tsfinder more verbose")
	rootCmd.PersistentFlags().BoolVarP(&color, "color", "c", false, "make tsfinder print with color")
	rootCmd.Execute()
}

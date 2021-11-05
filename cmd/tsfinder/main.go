// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/ariary/TrojanSourceFinder/pkg/bidirectional"
	"github.com/ariary/TrojanSourceFinder/pkg/homoglyph"
	"github.com/spf13/cobra"
)

func main() {
	var recursive bool
	var verbose bool
	var color bool

	//CMD FIND HOMOGLYPH
	var cmdSHomoglyph = &cobra.Command{
		Use:   "homoglyph [filename]",
		Short: "Detect homoglyph in file",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			homoglyph.Scan(name, recursive, verbose, color)

		},
	}

	var rootCmd = &cobra.Command{
		Use:   "tsfinder [filename]",
		Short: "Detect Trojan Source Vulnerability in your file",
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

	rootCmd.AddCommand(cmdSHomoglyph)
	rootCmd.Execute()
}

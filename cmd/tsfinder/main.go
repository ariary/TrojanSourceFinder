// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/ariary/TrojanSourceFinder/pkg/bidirectional"
	"github.com/ariary/TrojanSourceFinder/pkg/config"
	"github.com/ariary/TrojanSourceFinder/pkg/homoglyph"
	"github.com/spf13/cobra"
)

func main() {
	var verbose bool
	var color bool
	var onlyText bool
	var exclude string

	var sibling []string
	//CMD FIND HOMOGLYPH
	var cmdSHomoglyph = &cobra.Command{
		Use:   "homoglyph [path]",
		Short: "Detect homoglyph in file or directory",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			cfg := &config.Config{Verbose: verbose, Color: color, Sibling: &sibling, OnlyText: onlyText, ExcludelistFilename: exclude}

			homoglyph.Scan(path, cfg)
		},
	}
	//flag handling
	cmdSHomoglyph.PersistentFlags().StringSliceVarP(&sibling, "sibling", "s", nil, "(experimental) scan all files defined in scope to find the sibling  (ie word with the same skeleton) of homoglyphes found")

	//TSFINDER CMD
	var rootCmd = &cobra.Command{
		Use:   "tsfinder [path]",
		Short: "Detect Trojan Source Vulnerability in your file or directory",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			name := args[0]
			cfg := &config.Config{Verbose: verbose, Color: color, OnlyText: onlyText, ExcludelistFilename: exclude}

			bidirectional.Scan(name, cfg)
		},
	}

	//flag handling
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "make tsfinder more verbose")
	rootCmd.PersistentFlags().BoolVarP(&color, "color", "c", false, "make tsfinder print with color")
	rootCmd.PersistentFlags().BoolVarP(&onlyText, "text-file", "t", false, "make tsfinder scan only on 'human readable' file (ie  looks like correct UTF-8). Add verbosity (-v) to see which files has been skipped. This could help to rule out false positives")
	rootCmd.PersistentFlags().StringVarP(&exclude, "exclude", "e", "", "specifies a file containing a list of files not to be scanned. One file per line.")

	rootCmd.AddCommand(cmdSHomoglyph)
	rootCmd.Execute()
}

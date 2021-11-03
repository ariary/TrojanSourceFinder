// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"reflect"
	"unsafe"

	"github.com/ariary/TrojanSourceFinder/pkg/bidirectionnal"
	"github.com/spf13/cobra"
	"golang.org/x/text/unicode/bidi"
)

func Example_is() {

	// // constant with mixed type runes
	const mixed = "/*‮ } ⁦if (isAdmin)⁩ ⁦ begin admins only */"
	normal := "normale"
	var lr string
	lr = `العاشر ليونيكود (Unicode Conference)، الذي سيعقد في 10-12 آذار 1997 مبدينة`
	//mixed := "toto"
	// for _, c := range mixed {
	// 	if bidirectionnal.IsBidirectionalAlgorithm(c) {
	// 		fmt.Printf("%q: Is BDA\n", c)
	// 	}
	// }
	var p bidi.Paragraph
	p.SetBytes([]byte(lr))
	ord, err := p.Order()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%+v", ord.Direction())
	fmt.Println()
	rs := reflect.ValueOf(ord)
	rf := rs.FieldByName("directions")
	fmt.Printf("%+v\n", rf.Len())
	rs2 := reflect.New(rs.Type()).Elem()
	rs2.Set(rs)
	rf = rs2.FieldByName("directions")
	rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()

	fmt.Printf("%+v\n", rs.Type())
	fmt.Printf("%+v\n", rf.Len())
	test := rf.Index(0)
	fmt.Printf("%+v", test)

	fmt.Println(bidirectionnal.ContainBidirectionnal([]byte(lr)))
	fmt.Println(bidirectionnal.ContainBidirectionnal([]byte(normal)))

	// rs2 := reflect.New(rs.Type()).Elem()
	// rs2.Set(rs)
	// rf = rs2.Field(0)
	// rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
}

func main() {
	//CMD SCAN
	var recursive bool

	var cmdScan = &cobra.Command{
		Use:   "scan [filename/folder name]",
		Short: "Scan for Unicode Bidirectionnal characters within files",
		Long:  `Scan for Unicode Bidirectionnal characters within files as defined by Unicode's bidirectional algorithm property; this could lead to Trojan Source vulnerability`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			name := args[0]

			if recursive {
				//scan all files
			} else {
				//scan 1 file
			}

		},
	}

	//flag handling
	cmdScan.PersistentFlags().BoolVarP(&recursive, "recursive", "r", false, "scan all the files in the specified folder")

	//CMD EXORCISED

	var cmdExorcised = &cobra.Command{
		Use:   "exorcised [filename]",
		Short: "Print file without unicode bidirectionnal Algorithm",
		Long:  `print file without following unicode bidirectionnal Algorithm. It also print Unicode Bidirectionnal Character code point within text to help visualize where the direction changes happen.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			file := args[0]

			//bidirectionnalExorcisedPrint(file)

		},
	}
	var rootCmd = &cobra.Command{
		Use:   "tsFinder",
		Short: "Detect Trojan Source Vulnerability in your code",
	}
	rootCmd.AddCommand(cmdScan)
	rootCmd.AddCommand(cmdExorcised)
	rootCmd.Execute()
	Example_is()
}

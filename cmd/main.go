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
	Example_is()
}
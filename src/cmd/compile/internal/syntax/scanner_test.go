// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syntax

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestScanner(t *testing.T) {
	src, err := ioutil.ReadFile("parser.go")
	if err != nil {
		t.Fatal(err)
	}
	var s scanner
	s.init(src)
	for {
		s.next()
		if s.tok == _EOF {
			break
		}
		switch s.tok {
		case _Name:
			fmt.Println(s.line, s.tok, "=>", s.lit)
		case _Operator:
			fmt.Println(s.line, s.tok, "=>", s.op, s.prec)
		default:
			fmt.Println(s.line, s.tok)
		}
	}
}
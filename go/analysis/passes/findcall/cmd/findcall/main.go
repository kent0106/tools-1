// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The findcall command runs the findcall analyzer.
package main

import (
	"github.com/kent0106/gotools/go/analysis/passes/findcall"
	"github.com/kent0106/gotools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(findcall.Analyzer) }

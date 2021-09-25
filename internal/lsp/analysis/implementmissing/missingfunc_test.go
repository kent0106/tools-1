// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package implementmissing_test

import (
	"testing"

	"github.com/kent0106/gotools/go/analysis/analysistest"
	"github.com/kent0106/gotools/internal/lsp/analysis/implementmissing"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, implementmissing.Analyzer, "missingfunction")
}

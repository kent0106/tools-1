// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmdtest

import (
	"testing"

	"github.com/kent0106/gotools/internal/lsp/protocol"
	"github.com/kent0106/gotools/internal/span"
)

func (r *runner) Symbols(t *testing.T, uri span.URI, expectedSymbols []protocol.DocumentSymbol) {
	filename := uri.Filename()
	got, _ := r.NormalizeGoplsCmd(t, "symbols", filename)
	expect := string(r.data.Golden("symbols", filename, func() ([]byte, error) {
		return []byte(got), nil
	}))
	if expect != got {
		t.Errorf("symbols failed for %s expected:\n%s\ngot:\n%s", filename, expect, got)
	}
}

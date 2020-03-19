// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"context"
	"sync/atomic"
	"unsafe"
)

// Exporter is a function that handles events.
// It may return a modified context and event.
type Exporter func(context.Context, Event) (context.Context, Event)

var (
	exporter unsafe.Pointer
)

// SetExporter sets the global exporter function that handles all events.
// The exporter is called synchronously from the event call site, so it should
// return quickly so as not to hold up user code.
func SetExporter(e Exporter) {
	p := unsafe.Pointer(&e)
	if e == nil {
		// &e is always valid, and so p is always valid, but for the early abort
		// of ProcessEvent to be efficient it needs to make the nil check on the
		// pointer without having to dereference it, so we make the nil function
		// also a nil pointer
		p = nil
	}
	atomic.StorePointer(&exporter, p)
}

// ProcessEvent is called to deliver an event to the global exporter.
func ProcessEvent(ctx context.Context, ev Event) (context.Context, Event) {
	exporterPtr := (*Exporter)(atomic.LoadPointer(&exporter))
	if exporterPtr == nil {
		return ctx, ev
	}
	// and now also hand the event of to the current exporter
	return (*exporterPtr)(ctx, ev)
}

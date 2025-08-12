package perplexity

import (
	"context"
)

//
//-- Types
//

type optsKey struct{}

//
//-- With functions
//

// Adds an options struct to a context.
func WithOptions(ctx context.Context, opts *DLOpts) context.Context {
	return context.WithValue(ctx, optsKey{}, opts)
}

//
//-- From functions
//

// Pulls an options struct out of a context, falling back to default options if none is present.
func OptsFromCtx(ctx context.Context) *DLOpts {
	opts, ok := ctx.Value(optsKey{}).(*DLOpts)
	if !ok || opts == nil {
		o := DefaultDLOpts()
		return &o
	}
	return opts
}

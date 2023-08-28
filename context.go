package slogctx

import (
	"context"
	"log/slog"
	"sync"
)

type contextKey int

const attrsKey contextKey = iota

type attrs struct {
	attrs []slog.Attr
	mutex *sync.RWMutex
}

func WithAttrs(ctx context.Context, attrsToAdd ...slog.Attr) context.Context {
	ctxAttrs, ok := ctx.Value(attrsKey).(*attrs)
	if !ok {
		ctxAttrs = &attrs{
			mutex: &sync.RWMutex{},
		}
		ctx = context.WithValue(ctx, attrsKey, ctxAttrs)
	}
	ctxAttrs.mutex.Lock()
	ctxAttrs.attrs = append(ctxAttrs.attrs, attrsToAdd...)
	ctxAttrs.mutex.Unlock()
	return ctx
}

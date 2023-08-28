package slogctx

import (
	"context"
	"log/slog"
)

type ContextAttrFunc func(context.Context) []slog.Attr

func SetupHandler(handler slog.Handler, contextAttrFuncs ...ContextAttrFunc) slog.Handler {
	return &contextHandler{
		Handler:   handler,
		attrFuncs: contextAttrFuncs,
	}
}

type contextHandler struct {
	slog.Handler
	attrFuncs []ContextAttrFunc
}

func (handler *contextHandler) Handle(ctx context.Context, record slog.Record) error {
	var resolvedAttrsFromFuncs []slog.Attr
	for _, attrFunc := range handler.attrFuncs {
		resolvedAttrsFromFuncs = append(resolvedAttrsFromFuncs, attrFunc(ctx)...)
	}
	record.AddAttrs(resolvedAttrsFromFuncs...)

	ctxAttrs, ok := ctx.Value(attrsKey).(*attrs)
	if ok {
		ctxAttrs.mutex.RLock()
		record.AddAttrs(ctxAttrs.attrs...)
		ctxAttrs.mutex.RUnlock()
	}
	return handler.Handler.Handle(ctx, record)
}

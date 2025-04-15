package logger

import (
	"context"
	"errors"
	"log/slog"
)

type logCtx struct {
	Message string
}

type keyType int

const key = keyType(0)

func WithMessage(ctx context.Context, msg string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Message = msg
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{Message: msg})
}

type errorWithLogCtx struct {
	next error
	ctx  logCtx
}

func (e *errorWithLogCtx) Error() string {
	return e.next.Error()
}

func (e *errorWithLogCtx) Unwrap() error {
	return e.next
}

func WrapError(ctx context.Context, err error) error {
	c := logCtx{}
	if x, ok := ctx.Value(key).(logCtx); ok {
		c = x
	}

	return &errorWithLogCtx{
		next: err,
		ctx:  c,
	}
}

func ErrorCtx(ctx context.Context, err error) context.Context {
	myErr := new(errorWithLogCtx)
	if errors.As(err, &myErr) {
		return context.WithValue(ctx, key, myErr.ctx)
	}
	return ctx
}

func WrapErrorWithMessage(ctx context.Context, msg string, err error) error {
	return WrapError(WithMessage(ctx, msg), err)
}

func ExtractSlogAttrs(ctx context.Context) []slog.Attr {
	attr := make([]slog.Attr, 0, 8)

	if c, ok := ctx.Value(key).(logCtx); ok { // TODO if key is empty, will here panic?
		if c.Message != "" {
			attr = append(attr, slog.String("additional_message", c.Message))
		}
	}
	return attr
}

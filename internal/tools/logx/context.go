package logx

import "context"

type fieldsKey struct{}

// CtxWithFields returns a context carrying the given Fields.
func CtxWithFields(ctx context.Context, fields *Fields) context.Context {
	return context.WithValue(ctx, fieldsKey{}, fields)
}

// CtxFields returns the Fields stored in ctx, or an empty set.
func CtxFields(ctx context.Context) *Fields {
	if f, ok := ctx.Value(fieldsKey{}).(*Fields); ok {
		return f
	}
	return NewFields()
}

package context

import "context"

// SubjectCtxKey to be used with auth provider values
type subjectCtxKey struct {
}
type usernameCtxKey struct {
}

// subjectKey use to get and set auth provider values on context
var subjectKey = subjectCtxKey{}

// usernameKey use to get and set username on context
var usernameKey = usernameCtxKey{}

func Subject(ctx context.Context) (string, bool) {
	sub, ok := ctx.Value(subjectKey).(string)
	return sub, ok
}

func AddSubject(parent context.Context, subject string) context.Context {
	return context.WithValue(parent, subjectKey, subject)
}

func Username(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(usernameKey).(string)
	if username == "" {
		return "", false
	}
	return username, ok
}

func AddUsername(parent context.Context, username string) context.Context {
	return context.WithValue(parent, usernameKey, username)
}

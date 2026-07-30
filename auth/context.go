package auth

import "context"

type authenticatedSessionContextKey string

const authenticatedSessionContextKeyValue authenticatedSessionContextKey = "goc.auth.AuthenticatedSession"

// InjectAuthenticatedSession stores the authenticated session in the context.
func InjectAuthenticatedSession(ctx context.Context, session *AuthenticatedSession) context.Context {
	if session == nil {
		return ctx
	}
	return context.WithValue(ctx, authenticatedSessionContextKeyValue, session)
}

// ExtractAuthenticatedSession retrieves the authenticated session from the context.
func ExtractAuthenticatedSession(ctx context.Context) (*AuthenticatedSession, error) {
	session, ok := ctx.Value(authenticatedSessionContextKeyValue).(*AuthenticatedSession)
	if !ok || session == nil {
		return nil, ErrAuthenticatedSessionNotFound
	}
	return session, nil
}

// Session retrieves the authenticated session from the context.
func Session(ctx context.Context) (*AuthenticatedSession, bool) {
	session, err := ExtractAuthenticatedSession(ctx)
	return session, err == nil
}

// User retrieves the authenticated user from the context.
func User(ctx context.Context) (*AuthenticatedUser, bool) {
	session, ok := Session(ctx)
	if !ok {
		return nil, false
	}
	return &session.User, true
}

// UserID retrieves the authenticated user ID from the context.
func UserID(ctx context.Context) (string, bool) {
	user, ok := User(ctx)
	if !ok || user.ID == "" {
		return "", false
	}
	return user.ID, true
}

// UserInt64ID retrieves the authenticated user ID from the context as int64.
func UserInt64ID(ctx context.Context) (int64, bool) {
	user, ok := User(ctx)
	if !ok {
		return 0, false
	}
	id, err := user.Int64ID()
	if err != nil {
		return 0, false
	}
	return id, true
}

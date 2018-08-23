package web

// SessionAware is an action that injects the session into the context, it acquires a read lock on session.
func SessionAware(action Action) Action {
	return sessionMiddleware(action, nil)
}

// SessionRequired is an action that requires a session to be present
// or identified in some form on the request, and acquires a read lock on session.
func SessionRequired(action Action) Action {
	return sessionMiddleware(action, AuthManagerLoginRedirect)
}

// AuthManagerLoginRedirect is a redirect.
func AuthManagerLoginRedirect(ctx *Ctx) Result {
	return ctx.Auth().LoginRedirect(ctx)
}

// SessionMiddleware creates a custom session middleware.
func SessionMiddleware(notAuthorized Action) Middleware {
	return func(action Action) Action {
		return sessionMiddleware(action, notAuthorized)
	}
}

// SessionMiddleware returns a session middleware.
func sessionMiddleware(action, notAuthorized Action) Action {
	return func(ctx *Ctx) Result {
		session, err := ctx.Auth().VerifySession(ctx)
		if err != nil && !IsErrSessionInvalid(err) {
			return ctx.DefaultResultProvider().InternalError(err)
		}

		if session == nil {
			if notAuthorized != nil {
				return notAuthorized(ctx)
			}
			return action(ctx)
		}

		ctx.WithSession(session)
		return action(ctx)
	}
}

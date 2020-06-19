package web

// Middleware is a function that wraps a handler, it
// is called before the handler it wraps is called
type Middleware func(Handler) Handler

func wrapMiddleware(mw []Middleware, handler Handler) Handler {
	// loop backwards to make sure middleware is
	// executed as they are provided
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}
	return handler
}

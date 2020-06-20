// Package middleware provides primitives for configuring
// or reconfiguring HTTP handlers.
// they wrap web.Handler and are able to provide extra
// functionality that that are not neccessary to add into our HTTP handlers
// it can be used to enforce rules, collects metrics etc
package middleware

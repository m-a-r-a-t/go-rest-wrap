package middleware

import (
	"net/http"
)

type Middleware func(next http.HandlerFunc) http.HandlerFunc

type Chain struct {
	middlewares []Middleware
}

func New(dequeFuncs ...Middleware) Chain {
	return Chain{middlewares: dequeFuncs}
}

func (c Chain) Then(fn http.HandlerFunc) http.HandlerFunc {

	for i := range c.middlewares {
		fn = c.middlewares[len(c.middlewares)-1-i](fn)
	}

	return fn

}

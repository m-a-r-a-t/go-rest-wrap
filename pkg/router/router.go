package router

import (
	"net/http"
)

type Router struct {
	Mux         http.ServeMux
	handleFuncs map[string]*Funcs
}

type Funcs struct {
	get  func(w http.ResponseWriter, r *http.Request)
	post func(w http.ResponseWriter, r *http.Request)
}

func NewRouter() Router {
	return Router{
		Mux:         *http.NewServeMux(),
		handleFuncs: map[string]*Funcs{},
	}
}

func (r *Router) createHandleFunc(path string) {

	r.handleFuncs[path] = &Funcs{}
	r.Mux.HandleFunc(path, func(write http.ResponseWriter, request *http.Request) {

		switch true {
		case request.Method == "GET" && r.handleFuncs[path].get != nil:
			r.handleFuncs[path].get(write, request)
		case request.Method == "POST" && r.handleFuncs[path].post != nil:
			r.handleFuncs[path].post(write, request)
		default:
			http.NotFound(write, request)

		}

	})

}

func (r *Router) GET(path string, fn func(w http.ResponseWriter, r *http.Request)) {

	if _, ok := r.handleFuncs[path]; !ok {
		r.createHandleFunc(path)
	}

	r.handleFuncs[path].get = fn

}

func (r *Router) POST(path string, fn func(w http.ResponseWriter, r *http.Request)) {

	if _, ok := r.handleFuncs[path]; !ok {
		r.createHandleFunc(path)
	}

	r.handleFuncs[path].post = fn
}

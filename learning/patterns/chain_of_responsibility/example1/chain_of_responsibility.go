package example1

import "net/http"

type Middleware interface {
	handle(r http.Request) http.Response
}

type UserExistsMiddleware struct {
	next Middleware
}

func (u *UserExistsMiddleware) handle(r http.Request) http.Response {
	if r.URL.Query().Get("id") != "8" && r.URL.Query().Get("id") != "9" {
		return http.Response{Status: "400"}
	}
	return u.next.handle(r)
}

type SecurityMiddleware struct {
	next Middleware
}

func (s *SecurityMiddleware) handle(r http.Request) http.Response {
	if r.URL.Query().Get("id") != "9" {
		return http.Response{Status: "403"}
	}
	return s.next.handle(r)
}

type MainMiddleware struct {
	next Middleware
}

func (m *MainMiddleware) handle(r http.Request) http.Response {
	return http.Response{Status: "200"}
}

type Server struct{}

func (s *Server) handle(r http.Request) http.Response {
	middleware := &UserExistsMiddleware{
		next: &SecurityMiddleware{
			next: &MainMiddleware{next: nil},
		},
	}
	return middleware.handle(r)
}

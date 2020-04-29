package handlers

// package handlers imports http, and handles all communication between the
// client and server.

import (
	"net/http"
	"path"
	"strings"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	// TODO: import from github
	"../which"
)

func ListenAndServe(address string, handler *Root) error {
	return http.ListenAndServe(address, handler)
}

// == Root handler ================================

type Root struct {
	// child handlers
	StaticHandler *Static
	AuthHandler   *Auth
	//QsHandler *QsHandler
}

func NewRoot(userStore which.UserStore, sessionStore which.SessionStore) *Root {
	root := &Root{
		StaticHandler: &Static{},
		AuthHandler:   NewAuth(userStore, sessionStore),
	}
	return root
}

func (handler *Root) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	head, tail := shiftPath(req.URL.Path)
	if head == "qs" {
		req.URL.Path = tail
		//handler.QuestionsHandler.ServeHTTP(resp, req)
	} else if head == "auth" {
		req.URL.Path = tail
		handler.AuthHandler.ServeHTTP(resp, req)
	} else {
		req.URL.Path = head + tail
		handler.StaticHandler.ServeHTTP(resp, req)
	}
}

// == Static handler ================================
// serves files generated by Svelte (in /public)

type Static struct{}

func (handler *Static) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	http.ServeFile(resp, req, cleanPublicPath(req.URL.Path))
}

func cleanPublicPath(p string) string {
	cleaned := path.Clean("public/" + p)
	if cleaned != "public" && cleaned[:7] != "public/" {
		return "public/"
	}
	return cleaned
}

// == Shared helpers ================================

func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

func addCookie(resp http.ResponseWriter, name string, value string, expires time.Time) {
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expires,
		Path:    "/",
	}
	http.SetCookie(resp, &cookie)
}

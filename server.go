package cherry

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	// cherry version.
	VERSION = "0.1.0"
)

type Server struct {
	server *http.Server
	router *Router

	Port    int
	WebRoot string
	ResRoot string
	Name    string
}

func NewServer() *Server {
	s := &Server{}
	s.server = &http.Server{}
	s.router = NewRouter()
	s.Port = 0
	s.WebRoot = "WebRoot"
	s.ResRoot = "ResRoot"
	s.Name = fmt.Sprintf("cherry/%s", VERSION)
	return s
}

func (this *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// set headers
	w.Header().Set("Server", this.Name)
	// get url path
	urlPath := path.Clean(r.URL.Path)
	// static file request
	if len(path.Ext(urlPath)) > 1 && r.Method == http.MethodGet {
		file := this.WebRoot + urlPath
		info, err := os.Lstat(file)
		if err == nil && !info.IsDir() {
			http.ServeFile(w, r, file)
		} else {
			HttpErr(w, r, 404)
		}
		return
	}
	// service request
	route := this.router.FindRoute(r.Method, urlPath)
	if route != nil {
		values := r.URL.Query()
		a := strings.Split(route.Pattern, "/")
		b := strings.Split(urlPath, "/")
		for i, s := range a {
			if strings.HasPrefix(s, ":") {
				values.Add(a[i], b[i])
			}
		}
		r.URL.RawQuery = values.Encode()
		route.Business.Handle(w, r)
	} else {
		HttpErr(w, r, 404)
	}
}

func (this *Server) Serve() {
	// init businesses
	for _, r := range this.router.Routes {
		r.Business.Init(this.ResRoot)
	}
	// server configuration
	if this.Port == 0 {
		this.Port = 80
	}
	addr := fmt.Sprintf(":%d", this.Port)
	this.server.Addr = addr
	this.server.Handler = this
	// listen and serve
	err := this.server.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (this *Server) ServeTLS(certFile string, keyFile string) {
	// init businesses
	for _, r := range this.router.Routes {
		r.Business.Init(this.ResRoot)
	}
	// server configuration
	if this.Port == 0 {
		this.Port = 443
	}
	addr := fmt.Sprintf(":%d", this.Port)
	this.server.Addr = addr
	this.server.Handler = this
	// listen and serve
	err := this.server.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (this *Server) RouteGet(pattern string, b Business) {
	this.router.AddRoute(http.MethodGet, pattern, b)
}

func (this *Server) RoutePost(pattern string, b Business) {
	this.router.AddRoute(http.MethodGet, pattern, b)
}

func (this *Server) RoutePut(pattern string, b Business) {
	this.router.AddRoute(http.MethodPut, pattern, b)
}

func (this *Server) RouteDelete(pattern string, b Business) {
	this.router.AddRoute(http.MethodDelete, pattern, b)
}

func (this *Server) RoutePatch(pattern string, b Business) {
	this.router.AddRoute(http.MethodPatch, pattern, b)
}

func (this *Server) RouteHead(pattern string, b Business) {
	this.router.AddRoute(http.MethodHead, pattern, b)
}

package server

import (
	"net/http"
	"fmt"
	"strconv"
	"io/ioutil"
	"path/filepath"
)

type Server struct{
	Cache map[string] string
	On404 func(http.ResponseWriter,*http.Request)
}

func (s *Server) AddRoute(path string, callback func(http.ResponseWriter,*http.Request)) *Server {
	http.HandleFunc(path, callback)
	return s
}

func (s *Server) AddStaticPage(routePath string, pagePath string) *Server {
	absPath, _:= filepath.Abs(pagePath)
	body, err := ioutil.ReadFile(absPath)
	if err != nil{
		s.Cache[routePath] = "There was an error reading the file :C.\n\tError: " + fmt.Sprint(err)
		fmt.Println(err)
	} else {
		s.Cache[routePath] = string(body)
	}
	http.HandleFunc(routePath, s.LoadCachedPage)
	return s
}

func (s *Server) AddStaticPageFunc(pagePath string) func(http.ResponseWriter,*http.Request) {
	absPath, _:= filepath.Abs(pagePath)
	body, err := ioutil.ReadFile(absPath)
	if err == nil{
		return func(res http.ResponseWriter, req *http.Request){
			fmt.Fprintf(res, string(body))
		}
	} else {
		return func(res http.ResponseWriter, req *http.Request){
			fmt.Fprintf(res, "There was an error reading the file :C.\n\tError: " + fmt.Sprint(err))
		}
	}
}

func (s *Server) Start(port int) *Server {
	http.ListenAndServe(":"+strconv.Itoa(port),nil)
	return s
}

func (s *Server) LoadCachedPage(res http.ResponseWriter, req *http.Request){
	val, ok := s.Cache[req.URL.Path]
	if ok {
		fmt.Fprintf(res,val)
	} else {
		s.On404(res,req)
	}
}

func (s *Server) AddStaticFileserver(routePath string, dirPath string) *Server {
	absPath, _:= filepath.Abs(dirPath)
	fs := http.FileServer(http.Dir(absPath))
	http.Handle(routePath, http.StripPrefix(routePath, fs))
	return s
}

func NewServer() *Server{
	s := Server{}
	s.Cache = make(map[string] string)
	s.On404 = func(res http.ResponseWriter, req *http.Request){fmt.Fprintf(res,"404")}
	return &s
}
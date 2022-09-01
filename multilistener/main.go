package main

import (
	"log"
	"net/http"
)

type route struct {
	addr    string
	path    string
	handler http.HandlerFunc
}

var table = []route{
	{":8001", "/", handlerRoot},
	{":8001", "/path1", handlerPath1},
	{":8002", "/", handlerRoot},
	{":8002", "/path2", handlerPath2},
}

type server struct {
	server *http.Server
	mux    *http.ServeMux
}

func main() {

	serverTab := map[string]server{}

	for _, r := range table {
		s, found := serverTab[r.addr]
		if !found {
			mux := http.NewServeMux()
			s = server{
				server: &http.Server{
					Addr:    r.addr,
					Handler: mux,
				},
				mux: mux,
			}
			serverTab[r.addr] = s
		}
		s.mux.HandleFunc(r.path, r.handler)
		log.Printf("registered %v on port %s path %s", r.handler, r.addr, r.path)
	}

	for addr, s := range serverTab {
		go listenAndServe(s.server, addr)
	}

	<-chan struct{}(nil)
}

func listenAndServe(s *http.Server, addr string) {
	log.Printf("listening on port %s", addr)
	err := s.ListenAndServe()
	log.Printf("listening on port %s: %v", addr, err)
}

func handlerRoot(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not found", http.StatusNotFound)
}

func handlerPath1(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "path1", http.StatusOK)
}

func handlerPath2(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "path2", http.StatusOK)
}

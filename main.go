package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

func main() {
	dataDir := flag.String("data-dir", "/run/dump1090-mutability", "data dir")
	htmlDir := flag.String("html-dir", "/usr/share/dump1090-mutability/html", "html dir")
	port := flag.Int("port", 8081, "port")

	mux := http.NewServeMux()

	dataServer := http.FileServer(http.Dir(*dataDir))
	htmlServer := http.FileServer(http.Dir(*htmlDir))

	mux.Handle("/dump1090/data/", http.StripPrefix("/dump1090/data/", dataServer))
	mux.Handle("/dump1090", http.RedirectHandler("/dump1090/", http.StatusMovedPermanently))
	mux.HandleFunc("/dump1090/", func(writer http.ResponseWriter, request *http.Request) {

		parse, err := url.Parse(request.RequestURI)
		if err != nil {
			log.Printf("Failed to parse request URI: %s\n", request.RequestURI)
			writer.WriteHeader(http.StatusInternalServerError)
		}

		matched, err := regexp.MatchString("^/dump1090/$", parse.Path)
		if err != nil {
			log.Printf("Failed to match regex for path: %s", parse.Path)
			writer.WriteHeader(http.StatusInternalServerError)
		}
		if matched {
			http.RedirectHandler("/dump1090/gmap.html", http.StatusMovedPermanently).ServeHTTP(writer, request)
		} else {
			http.StripPrefix("/dump1090/", htmlServer).ServeHTTP(writer, request)
		}
	})

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

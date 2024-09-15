package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	urlshort "github.com/rdustin/gophercises-learning/url-shortener"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v3",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := []byte(`
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`)
	json := []byte(`
	[{"path": "/b", "url": "https://google.com"}]
	`)
	jsonPath := flag.String("j", "", "Path to json file denoting path and url combinations. (default empty)")
	var err error
	yamlPath := flag.String("y", "", "Path to yaml file denoting path and url combinations. (default empty)")
	flag.Parse()
	if *yamlPath != "" {
		yaml, err = os.ReadFile(*yamlPath)
		if err != nil {
			panic(err)
		}
	}
	yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
	if err != nil {
		panic(err)
	}

	// Build JsonHandler using the YAMLHandler as the fallback
	if *jsonPath != "" {
		json, err = os.ReadFile(*jsonPath)
		if err != nil {
			panic(err)
		}
	}
	jsonHandler, err := urlshort.JsonHandler(json, yamlHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

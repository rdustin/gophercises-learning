package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
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
	_, err = setupDb()
	if err != nil {
		fmt.Print(err)
	}
	dbHandler, err := urlshort.DbHandler(jsonHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", dbHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func setupDb() (*bolt.DB, error) {
	db, err := bolt.Open("paths.db", 0600, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	initDbData(db)
	return db, nil
}

func initDbData(db *bolt.DB) {
	data := []byte(`[{"path": "/z", "url": "https://google.com"}]`)
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("DB"))
		if err != nil {
			return fmt.Errorf("could not create root bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		err = tx.Bucket([]byte("DB")).Put([]byte("paths"), data)
		if err != nil {
			return fmt.Errorf("could not set data")
		}
		return nil
	})
}

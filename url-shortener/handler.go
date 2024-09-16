package urlshort

import (
	"encoding/json"
	"net/http"

	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		found := ""
		for path, url := range pathsToUrls {
			if r.URL.Path == path {
				found = url
			}
		}
		if found != "" {
			http.Redirect(w, r, found, http.StatusTemporaryRedirect)
			return
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yaml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func JsonHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJson, err := parseJson(json)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedJson)
	return MapHandler(pathMap, fallback), nil
}

func DbHandler(fallback http.Handler) (http.HandlerFunc, error) {
	// don't like doing it this way, but for some reason the db is closing before it gets here if I pass it in
	db, err := bolt.Open("paths.db", 0600, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var data []byte
	err = db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte("DB")).Get([]byte("paths"))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return JsonHandler(data, fallback)
}

func parseYAML(yml []byte) ([]pathItem, error) {
	parsedValues := []pathItem{}
	err := yaml.Unmarshal(yml, &parsedValues)
	if err != nil {
		return nil, err
	}
	return parsedValues, nil
}

func parseJson(jsn []byte) ([]pathItem, error) {
	parsedValues := []pathItem{}
	err := json.Unmarshal(jsn, &parsedValues)
	if err != nil {
		return nil, err
	}
	return parsedValues, nil
}

func buildMap(pathItems []pathItem) map[string]string {
	pathMap := map[string]string{}
	for _, pi := range pathItems {
		pathMap[pi.Path] = pi.Url
	}
	return pathMap
}

type pathItem struct {
	Path string
	Url  string
}

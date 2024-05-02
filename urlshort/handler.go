package urlshort

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	var f http.HandlerFunc
	f = func(w http.ResponseWriter, r *http.Request) {
		// See if the request path matches any of the predefined paths in map
		path := r.URL.Path
		log.Printf("Searching for path %s\n", path)
		if redirect_url, ok := pathsToUrls[path]; ok {
			// Found match, redirect
			http.Redirect(w, r, redirect_url, http.StatusSeeOther)

		} else {
			// No match found
			fallback.ServeHTTP(w, r)
		}
	}
	return f
}

type mapping struct {
	Path string
	Url  string
}

func buildMap(mappings []mapping) map[string]string {
	m := make(map[string]string)
	for _, mapping := range mappings {
		m[mapping.Path] = mapping.Url
	}
	return m
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
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	paths := make([]mapping, 1)
	err := yaml.Unmarshal(yml, &paths)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(paths)
	return MapHandler(pathsToUrls, fallback), nil
}

func JSONHandler(json_data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	paths := make([]mapping, 1)
	err := json.Unmarshal(json_data, &paths)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(paths)
	return MapHandler(pathsToUrls, fallback), nil
}
